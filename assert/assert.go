// Package assert provides basic assertions for unit tests.
//
// Example Usage
//
// The following examples show how to use this library:
//    import (
//      "testing"
//      . "v2ray.com/ext/assert"
//    )
//
//    func TestSomething(t *testing.T) {
//      assert := With(t)
//
//      a := 10
//      b := 1
//      assert(a, Equals, b)
//    }
//
package assert

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

var (
	equals  = make(map[reflect.Type]*Matcher)
	less    = make(map[reflect.Type]*Matcher)
	greater = make(map[reflect.Type]*Matcher)
)

type Matcher struct {
	method reflect.Value
	verb   string
}

func zeroValueOrReal(v interface{}, t reflect.Type) reflect.Value {
	if v == nil {
		return reflect.New(t).Elem()
	}
	return reflect.ValueOf(v)
}

func (op *Matcher) call(value interface{}, expectations []interface{}) bool {
	input := make([]reflect.Value, op.method.Type().NumIn())
	input[0] = zeroValueOrReal(value, op.method.Type().In(0))
	for i, v := range expectations {
		input[i+1] = zeroValueOrReal(v, op.method.Type().In(i+1))
	}
	ret := op.method.Call(input)
	return ret[0].Bool()
}

// Assert asserts the given value matches expectations.
type Assert func(value interface{}, op *Matcher, expectations ...interface{})

// With creates wrap the testing object into an assertion.
func With(t *testing.T) Assert {
	return func(value interface{}, op *Matcher, expectations ...interface{}) {
		if !op.call(value, expectations) {
			msgs := []interface{}{"Not true that (", value, ") ", op.verb, " ("}
			msgs = append(msgs, expectations...)
			msgs = append(msgs, ")")
			fmt.Println(decorate(fmt.Sprint(msgs...)))
			t.FailNow()
		}
	}
}

func storeToTable(m map[reflect.Type]*Matcher, f reflect.Value, op *Matcher) {
	t := f.Type().In(0)
	m[t] = op
}

func RegisterEqualsMatcher(f interface{}) {
	op := CreateMatcher(f, "equals to")
	storeToTable(equals, reflect.ValueOf(f), op)
}

func RegisterLessThanMatcher(f interface{}) {
	op := CreateMatcher(f, "less than")
	storeToTable(less, reflect.ValueOf(f), op)
}

func RegisterGreaterThanMatcher(f interface{}) {
	op := CreateMatcher(f, "greater than")
	storeToTable(greater, reflect.ValueOf(f), op)
}

func CreateMatcher(fv interface{}, verb string) *Matcher {
	f := reflect.ValueOf(fv)
	if f.Kind() != reflect.Func {
		panic("Operator is not a function.")
	}

	return &Matcher{
		method: f,
		verb:   verb,
	}
}

func callInternal(m map[reflect.Type]*Matcher, v interface{}, exp interface{}) bool {
	vt := reflect.TypeOf(v)
	op, found := m[vt]
	if !found {
		panic(fmt.Sprint("Type (", vt, ") not registered."))
	}
	return op.call(v, []interface{}{exp})
}

// Equals is a Matcher that expects two given values are equal to each other.
var Equals = CreateMatcher(func(v interface{}, exp interface{}) bool {
	vt := reflect.TypeOf(v)
	op, found := equals[vt]
	if found {
		return op.call(v, []interface{}{exp})
	}
	return v == exp
}, "equals to")

// NotEquals is a Matcher that expects two given values are not equal to each other.
var NotEquals = CreateMatcher(func(v interface{}, exp interface{}) bool {
	vt := reflect.TypeOf(v)
	op, found := equals[vt]
	if found {
		return !op.call(v, []interface{}{exp})
	}
	return v != exp
}, "not equals to")

// LessThan expects the given value is less than the expectation.
var LessThan = CreateMatcher(func(v interface{}, exp interface{}) bool {
	return callInternal(less, v, exp)
}, "less than")

var AtMost = CreateMatcher(func(v interface{}, exp interface{}) bool {
	return !callInternal(greater, v, exp)
}, "less than or equals to")

var GreaterThan = CreateMatcher(func(v interface{}, exp interface{}) bool {
	return callInternal(greater, v, exp)
}, "less than")

var AtLeast = CreateMatcher(func(v interface{}, exp interface{}) bool {
	return !callInternal(less, v, exp)
}, "less than")

var IsNegative = CreateMatcher(func(v interface{}, exp interface{}) bool {
	return callInternal(less, v, 0)
}, "is negative")

var IsPositive = CreateMatcher(func(v interface{}, exp interface{}) bool {
	return callInternal(greater, v, 0)
}, "is positive")

var IsNil = CreateMatcher(func(v interface{}) bool {
	return v == nil || reflect.ValueOf(v).IsNil()
}, "is nil")

var IsNotNil = CreateMatcher(func(v interface{}) bool {
	return v != nil
}, "is not nil")

var IsTrue = CreateMatcher(func(v bool) bool {
	return v
}, "is true")

var IsFalse = CreateMatcher(func(v bool) bool {
	return !v
}, "is false")

var IsEmpty = CreateMatcher(func(v interface{}) bool {
	return reflect.ValueOf(v).Len() == 0
}, "is empty")

var Panics = CreateMatcher(func(v interface{}) (ret bool) {
	defer func() {
		if x := recover(); x != nil {
			ret = true
		}
	}()
	if vf, ok := v.(func()); ok {
		vf()
	}
	return false
}, "panics")

var Implements = CreateMatcher(func(v interface{}, exp interface{}) bool {
	return reflect.TypeOf(v).Implements(reflect.TypeOf(exp).Elem())
}, "implements")

func Not(op *Matcher) *Matcher {
	return &Matcher{
		method: reflect.MakeFunc(op.method.Type(), func(v []reflect.Value) []reflect.Value {
			return []reflect.Value{reflect.ValueOf(!op.method.Call(v)[0].Bool())}
		}),
		verb: "not " + op.verb,
	}
}

func init() {
	RegisterEqualsMatcher(func(v, exp bool) bool {
		return v == exp
	})

	RegisterEqualsMatcher(func(v, exp byte) bool {
		return v == exp
	})

	RegisterEqualsMatcher(func(v, exp int8) bool {
		return v == exp
	})

	RegisterEqualsMatcher(func(v, exp uint8) bool {
		return v == exp
	})

	RegisterEqualsMatcher(func(v, exp int16) bool {
		return v == exp
	})

	RegisterEqualsMatcher(func(v, exp uint16) bool {
		return v == exp
	})

	RegisterEqualsMatcher(func(v int, exp int) bool {
		return v == exp
	})

	RegisterEqualsMatcher(func(v uint, exp uint) bool {
		return v == exp
	})

	RegisterEqualsMatcher(func(v, exp int32) bool {
		return v == exp
	})

	RegisterEqualsMatcher(func(v, exp uint32) bool {
		return v == exp
	})

	RegisterEqualsMatcher(func(v, exp int64) bool {
		return v == exp
	})

	RegisterEqualsMatcher(func(v, exp uint64) bool {
		return v == exp
	})

	RegisterLessThanMatcher(func(v, exp byte) bool {
		return v < exp
	})

	RegisterLessThanMatcher(func(v, exp int8) bool {
		return v < exp
	})

	RegisterLessThanMatcher(func(v, exp uint8) bool {
		return v < exp
	})

	RegisterLessThanMatcher(func(v, exp int16) bool {
		return v < exp
	})

	RegisterLessThanMatcher(func(v, exp uint16) bool {
		return v < exp
	})

	RegisterLessThanMatcher(func(v int, exp int) bool {
		return v < exp
	})

	RegisterLessThanMatcher(func(v, exp int32) bool {
		return v < exp
	})

	RegisterLessThanMatcher(func(v, exp uint32) bool {
		return v < exp
	})

	RegisterLessThanMatcher(func(v, exp int64) bool {
		return v < exp
	})

	RegisterLessThanMatcher(func(v, exp uint64) bool {
		return v < exp
	})

	RegisterGreaterThanMatcher(func(v, exp byte) bool {
		return v > exp
	})

	RegisterGreaterThanMatcher(func(v, exp int8) bool {
		return v > exp
	})

	RegisterGreaterThanMatcher(func(v, exp uint8) bool {
		return v > exp
	})

	RegisterGreaterThanMatcher(func(v, exp int16) bool {
		return v > exp
	})

	RegisterGreaterThanMatcher(func(v, exp uint16) bool {
		return v > exp
	})

	RegisterGreaterThanMatcher(func(v int, exp int) bool {
		return v > exp
	})

	RegisterGreaterThanMatcher(func(v, exp int32) bool {
		return v > exp
	})

	RegisterGreaterThanMatcher(func(v, exp uint32) bool {
		return v > exp
	})

	RegisterGreaterThanMatcher(func(v, exp int64) bool {
		return v > exp
	})

	RegisterGreaterThanMatcher(func(v, exp uint64) bool {
		return v > exp
	})

	RegisterEqualsMatcher(func(v, exp string) bool {
		return v == exp
	})

	RegisterEqualsMatcher(func(v, exp []byte) bool {
		if len(v) != len(exp) {
			return false
		}
		for i, vv := range v {
			if vv != exp[i] {
				return false
			}
		}

		return true
	})

	RegisterEqualsMatcher(func(v, exp []string) bool {
		if len(v) != len(exp) {
			return false
		}
		for i, vv := range v {
			if vv != exp[i] {
				return false
			}
		}

		return true
	})
}

func getCaller() (string, int) {
	stackLevel := 1
	for {
		_, file, line, ok := runtime.Caller(stackLevel)
		if strings.Contains(file, "assert") {
			stackLevel++
		} else {
			if ok {
				// Truncate file name at last file name separator.
				if index := strings.LastIndex(file, "/"); index >= 0 {
					file = file[index+1:]
				} else if index = strings.LastIndex(file, "\\"); index >= 0 {
					file = file[index+1:]
				}
			} else {
				file = "???"
				line = 1
			}
			return file, line
		}
	}
}

// decorate prefixes the string with the file and line of the call site
// and inserts the final newline if needed and indentation tabs for formatting.
func decorate(s string) string {
	file, line := getCaller()
	buf := new(bytes.Buffer)
	// Every line is indented at least one tab.
	buf.WriteString("  ")
	fmt.Fprintf(buf, "%s:%d: ", file, line)
	lines := strings.Split(s, "\n")
	if l := len(lines); l > 1 && lines[l-1] == "" {
		lines = lines[:l-1]
	}
	for i, line := range lines {
		if i > 0 {
			// Second and subsequent lines are indented an extra tab.
			buf.WriteString("\n\t\t")
		}
		buf.WriteString(line)
	}
	buf.WriteByte('\n')
	return buf.String()
}
