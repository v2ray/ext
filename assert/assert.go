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

type Assert func(value interface{}, op *Matcher, expectations ...interface{})

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

func RegisterEqualsMatcher(t reflect.Type, f reflect.Value) {
	op := CreateMatcher(t, f, 2, "equals to")
	equals[t] = op
}

func RegisterLessThanMatcher(t reflect.Type, f reflect.Value) {
	op := CreateMatcher(t, f, 2, "less than")
	less[t] = op
}

func RegisterGreaterThanMatcher(t reflect.Type, f reflect.Value) {
	op := CreateMatcher(t, f, 2, "greater than")
	greater[t] = op
}

func CreateMatcher(t reflect.Type, f reflect.Value, numInput int, verb string) *Matcher {
	if f.Kind() != reflect.Func {
		panic("Operator is not a function.")
	}

	if f.Type().NumIn() != numInput {
		panic(fmt.Sprint("Operator accepts ", f.Type().NumIn(), " parameters, but expect to be ", numInput))
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

var Equals = &Matcher{
	method: reflect.ValueOf(func(v interface{}, exp interface{}) bool {
		return callInternal(equals, v, exp)
	}),
	verb: "equals to",
}

var NotEquals = &Matcher{
	method: reflect.ValueOf(func(v interface{}, exp interface{}) bool {
		return !callInternal(equals, v, exp)
	}),
	verb: "not equals to",
}

var LessThan = &Matcher{
	method: reflect.ValueOf(func(v interface{}, exp interface{}) bool {
		return callInternal(less, v, exp)
	}),
	verb: "less than",
}

var LessThanOrEqualsTo = &Matcher{
	method: reflect.ValueOf(func(v interface{}, exp interface{}) bool {
		return !callInternal(greater, v, exp)
	}),
	verb: "less than or equals to",
}

var GreaterThan = &Matcher{
	method: reflect.ValueOf(func(v interface{}, exp interface{}) bool {
		return callInternal(greater, v, exp)
	}),
	verb: "less than",
}

var GreaterThanOrEqualsTo = &Matcher{
	method: reflect.ValueOf(func(v interface{}, exp interface{}) bool {
		return !callInternal(less, v, exp)
	}),
	verb: "less than",
}

var IsNegative = &Matcher{
	method: reflect.ValueOf(func(v interface{}, exp interface{}) bool {
		return callInternal(less, v, 0)
	}),
	verb: "is negative",
}

var IsPositive = &Matcher{
	method: reflect.ValueOf(func(v interface{}, exp interface{}) bool {
		return callInternal(greater, v, 0)
	}),
	verb: "is positive",
}

var IsNil = CreateMatcher(reflect.TypeOf(interface{}(nil)), reflect.ValueOf(func(v interface{}) bool {
	return v == nil
}), 1, "is nil")

var IsNotNil = CreateMatcher(reflect.TypeOf(interface{}(nil)), reflect.ValueOf(func(v interface{}) bool {
	return v != nil
}), 1, "is not nil")

var IsTrue = CreateMatcher(reflect.TypeOf(true), reflect.ValueOf(func(v bool) bool {
	return v
}), 1, "is true")

var IsFalse = CreateMatcher(reflect.TypeOf(true), reflect.ValueOf(func(v bool) bool {
	return !v
}), 1, "is false")

func Not(op *Matcher) *Matcher {
	return &Matcher{
		method: reflect.MakeFunc(op.method.Type(), func(v []reflect.Value) []reflect.Value {
			return []reflect.Value{reflect.ValueOf(!op.method.Call(v)[0].Bool())}
		}),
		verb: "not " + op.verb,
	}
}

func init() {
	RegisterEqualsMatcher(reflect.TypeOf(true), reflect.ValueOf(func(v, exp bool) bool {
		return v == exp
	}))

	RegisterEqualsMatcher(reflect.TypeOf(byte(0)), reflect.ValueOf(func(v, exp byte) bool {
		return v == exp
	}))

	RegisterEqualsMatcher(reflect.TypeOf(int8(0)), reflect.ValueOf(func(v, exp int8) bool {
		return v == exp
	}))

	RegisterEqualsMatcher(reflect.TypeOf(uint8(0)), reflect.ValueOf(func(v, exp uint8) bool {
		return v == exp
	}))

	RegisterEqualsMatcher(reflect.TypeOf(int16(0)), reflect.ValueOf(func(v, exp int16) bool {
		return v == exp
	}))

	RegisterEqualsMatcher(reflect.TypeOf(uint16(0)), reflect.ValueOf(func(v, exp uint16) bool {
		return v == exp
	}))

	RegisterEqualsMatcher(reflect.TypeOf(int(0)), reflect.ValueOf(func(v int, exp int) bool {
		return v == exp
	}))

	RegisterEqualsMatcher(reflect.TypeOf(uint(0)), reflect.ValueOf(func(v uint, exp uint) bool {
		return v == exp
	}))

	RegisterEqualsMatcher(reflect.TypeOf(int32(0)), reflect.ValueOf(func(v, exp int32) bool {
		return v == exp
	}))

	RegisterEqualsMatcher(reflect.TypeOf(uint32(0)), reflect.ValueOf(func(v, exp uint32) bool {
		return v == exp
	}))

	RegisterEqualsMatcher(reflect.TypeOf(int64(0)), reflect.ValueOf(func(v, exp int64) bool {
		return v == exp
	}))

	RegisterEqualsMatcher(reflect.TypeOf(uint64(0)), reflect.ValueOf(func(v, exp uint64) bool {
		return v == exp
	}))

	RegisterLessThanMatcher(reflect.TypeOf(byte(0)), reflect.ValueOf(func(v, exp byte) bool {
		return v < exp
	}))

	RegisterLessThanMatcher(reflect.TypeOf(int8(0)), reflect.ValueOf(func(v, exp int8) bool {
		return v < exp
	}))

	RegisterLessThanMatcher(reflect.TypeOf(uint8(0)), reflect.ValueOf(func(v, exp uint8) bool {
		return v < exp
	}))

	RegisterLessThanMatcher(reflect.TypeOf(int16(0)), reflect.ValueOf(func(v, exp int16) bool {
		return v < exp
	}))

	RegisterLessThanMatcher(reflect.TypeOf(uint16(0)), reflect.ValueOf(func(v, exp uint16) bool {
		return v < exp
	}))

	RegisterLessThanMatcher(reflect.TypeOf(int(0)), reflect.ValueOf(func(v int, exp int) bool {
		return v < exp
	}))

	RegisterLessThanMatcher(reflect.TypeOf(int32(0)), reflect.ValueOf(func(v, exp int32) bool {
		return v < exp
	}))

	RegisterLessThanMatcher(reflect.TypeOf(uint32(0)), reflect.ValueOf(func(v, exp uint32) bool {
		return v < exp
	}))

	RegisterLessThanMatcher(reflect.TypeOf(int64(0)), reflect.ValueOf(func(v, exp int64) bool {
		return v < exp
	}))

	RegisterLessThanMatcher(reflect.TypeOf(uint64(0)), reflect.ValueOf(func(v, exp uint64) bool {
		return v < exp
	}))

	RegisterGreaterThanMatcher(reflect.TypeOf(byte(0)), reflect.ValueOf(func(v, exp byte) bool {
		return v > exp
	}))

	RegisterGreaterThanMatcher(reflect.TypeOf(int8(0)), reflect.ValueOf(func(v, exp int8) bool {
		return v > exp
	}))

	RegisterGreaterThanMatcher(reflect.TypeOf(uint8(0)), reflect.ValueOf(func(v, exp uint8) bool {
		return v > exp
	}))

	RegisterGreaterThanMatcher(reflect.TypeOf(int16(0)), reflect.ValueOf(func(v, exp int16) bool {
		return v > exp
	}))

	RegisterGreaterThanMatcher(reflect.TypeOf(uint16(0)), reflect.ValueOf(func(v, exp uint16) bool {
		return v > exp
	}))

	RegisterGreaterThanMatcher(reflect.TypeOf(int(0)), reflect.ValueOf(func(v int, exp int) bool {
		return v > exp
	}))

	RegisterGreaterThanMatcher(reflect.TypeOf(int32(0)), reflect.ValueOf(func(v, exp int32) bool {
		return v > exp
	}))

	RegisterGreaterThanMatcher(reflect.TypeOf(uint32(0)), reflect.ValueOf(func(v, exp uint32) bool {
		return v > exp
	}))

	RegisterGreaterThanMatcher(reflect.TypeOf(int64(0)), reflect.ValueOf(func(v, exp int64) bool {
		return v > exp
	}))

	RegisterGreaterThanMatcher(reflect.TypeOf(uint64(0)), reflect.ValueOf(func(v, exp uint64) bool {
		return v > exp
	}))

	RegisterEqualsMatcher(reflect.TypeOf(""), reflect.ValueOf(func(v, exp string) bool {
		return v == exp
	}))

	RegisterEqualsMatcher(reflect.TypeOf([]byte(nil)), reflect.ValueOf(func(v, exp []byte) bool {
		if len(v) != len(exp) {
			return false
		}
		for i, vv := range v {
			if vv != exp[i] {
				return false
			}
		}

		return true
	}))

	RegisterEqualsMatcher(reflect.TypeOf([]string(nil)), reflect.ValueOf(func(v, exp []string) bool {
		if len(v) != len(exp) {
			return false
		}
		for i, vv := range v {
			if vv != exp[i] {
				return false
			}
		}

		return true
	}))
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
