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
	equals  = make(map[reflect.Type]*Operator)
	less    = make(map[reflect.Type]*Operator)
	greater = make(map[reflect.Type]*Operator)
)

type Operator struct {
	method reflect.Value
	verb   string
}

func zeroValueOrReal(v interface{}, t reflect.Type) reflect.Value {
	if v == nil {
		return reflect.New(t).Elem()
	}
	return reflect.ValueOf(v)
}

func (op *Operator) call(value interface{}, expectations []interface{}) bool {
	input := make([]reflect.Value, op.method.Type().NumIn())
	input[0] = zeroValueOrReal(value, op.method.Type().In(0))
	for i, v := range expectations {
		input[i+1] = zeroValueOrReal(v, op.method.Type().In(i+1))
	}
	ret := op.method.Call(input)
	return ret[0].Bool()
}

type Assert func(value interface{}, op *Operator, expectations ...interface{})

func With(t *testing.T) Assert {
	return func(value interface{}, op *Operator, expectations ...interface{}) {
		if !op.call(value, expectations) {
			fmt.Println(decorate(fmt.Sprint("Not true that ", value, " ", op.verb, " ", expectations)))
			t.FailNow()
		}
	}
}

func RegisterEqualsOperator(t reflect.Type, f reflect.Value) {
	op := CreateOperator(t, f, 2, "equals to")
	equals[t] = op
}

func RegisterLessThanOperator(t reflect.Type, f reflect.Value) {
	op := CreateOperator(t, f, 2, "less than")
	less[t] = op
}

func RegisterGreaterThanOperator(t reflect.Type, f reflect.Value) {
	op := CreateOperator(t, f, 2, "greater than")
	greater[t] = op
}

func CreateOperator(t reflect.Type, f reflect.Value, numInput int, verb string) *Operator {
	if f.Kind() != reflect.Func {
		panic("Operator is not a function.")
	}

	if f.Type().NumIn() != numInput {
		panic(fmt.Sprint("Operator accepts ", f.Type().NumIn(), " parameters, but expect to be ", numInput))
	}

	return &Operator{
		method: f,
		verb:   verb,
	}
}

func callInternal(m map[reflect.Type]*Operator, v interface{}, exp interface{}) bool {
	vt := reflect.TypeOf(v)
	op, found := m[vt]
	if !found {
		panic(fmt.Sprint("Type", vt, "not registered."))
	}
	return op.call(v, []interface{}{exp})
}

var Equals = &Operator{
	method: reflect.ValueOf(func(v interface{}, exp interface{}) bool {
		return callInternal(equals, v, exp)
	}),
	verb: "equals to",
}

var NotEquals = &Operator{
	method: reflect.ValueOf(func(v interface{}, exp interface{}) bool {
		return !callInternal(equals, v, exp)
	}),
	verb: "not equals to",
}

var LessThan = &Operator{
	method: reflect.ValueOf(func(v interface{}, exp interface{}) bool {
		return callInternal(less, v, exp)
	}),
	verb: "less than",
}

var LessThanOrEqualsTo = &Operator{
	method: reflect.ValueOf(func(v interface{}, exp interface{}) bool {
		return !callInternal(greater, v, exp)
	}),
	verb: "less than or equals to",
}

var GreaterThan = &Operator{
	method: reflect.ValueOf(func(v interface{}, exp interface{}) bool {
		return callInternal(greater, v, exp)
	}),
	verb: "less than",
}

var GreaterThanOrEqualsTo = &Operator{
	method: reflect.ValueOf(func(v interface{}, exp interface{}) bool {
		return !callInternal(less, v, exp)
	}),
	verb: "less than",
}

var IsNegative = &Operator{
	method: reflect.ValueOf(func(v interface{}, exp interface{}) bool {
		return callInternal(less, v, 0)
	}),
	verb: "is negative",
}

var IsPositive = &Operator{
	method: reflect.ValueOf(func(v interface{}, exp interface{}) bool {
		return callInternal(greater, v, 0)
	}),
	verb: "is positive",
}

var IsNil = CreateOperator(reflect.TypeOf(interface{}(nil)), reflect.ValueOf(func(v interface{}) bool {
	return v == nil
}), 1, "is nil")

var IsNotNil = CreateOperator(reflect.TypeOf(interface{}(nil)), reflect.ValueOf(func(v interface{}) bool {
	return v != nil
}), 1, "is not nil")

var IsTrue = CreateOperator(reflect.TypeOf(true), reflect.ValueOf(func(v bool) bool {
	return v
}), 1, "is true")

var IsFalse = CreateOperator(reflect.TypeOf(true), reflect.ValueOf(func(v bool) bool {
	return !v
}), 1, "is false")

func init() {
	RegisterEqualsOperator(reflect.TypeOf(true), reflect.ValueOf(func(v, exp bool) bool {
		return v == exp
	}))

	RegisterEqualsOperator(reflect.TypeOf(byte(0)), reflect.ValueOf(func(v, exp byte) bool {
		return v == exp
	}))

	RegisterEqualsOperator(reflect.TypeOf(int8(0)), reflect.ValueOf(func(v, exp int8) bool {
		return v == exp
	}))

	RegisterEqualsOperator(reflect.TypeOf(uint8(0)), reflect.ValueOf(func(v, exp uint8) bool {
		return v == exp
	}))

	RegisterEqualsOperator(reflect.TypeOf(int16(0)), reflect.ValueOf(func(v, exp int16) bool {
		return v == exp
	}))

	RegisterEqualsOperator(reflect.TypeOf(uint16(0)), reflect.ValueOf(func(v, exp uint16) bool {
		return v == exp
	}))

	RegisterEqualsOperator(reflect.TypeOf(int(0)), reflect.ValueOf(func(v int, exp int) bool {
		return v == exp
	}))

	RegisterEqualsOperator(reflect.TypeOf(uint(0)), reflect.ValueOf(func(v uint, exp uint) bool {
		return v == exp
	}))

	RegisterEqualsOperator(reflect.TypeOf(int32(0)), reflect.ValueOf(func(v, exp int32) bool {
		return v == exp
	}))

	RegisterEqualsOperator(reflect.TypeOf(uint32(0)), reflect.ValueOf(func(v, exp uint32) bool {
		return v == exp
	}))

	RegisterEqualsOperator(reflect.TypeOf(int64(0)), reflect.ValueOf(func(v, exp int64) bool {
		return v == exp
	}))

	RegisterEqualsOperator(reflect.TypeOf(uint64(0)), reflect.ValueOf(func(v, exp uint64) bool {
		return v == exp
	}))

	RegisterLessThanOperator(reflect.TypeOf(byte(0)), reflect.ValueOf(func(v, exp byte) bool {
		return v < exp
	}))

	RegisterLessThanOperator(reflect.TypeOf(int8(0)), reflect.ValueOf(func(v, exp int8) bool {
		return v < exp
	}))

	RegisterLessThanOperator(reflect.TypeOf(uint8(0)), reflect.ValueOf(func(v, exp uint8) bool {
		return v < exp
	}))

	RegisterLessThanOperator(reflect.TypeOf(int16(0)), reflect.ValueOf(func(v, exp int16) bool {
		return v < exp
	}))

	RegisterLessThanOperator(reflect.TypeOf(uint16(0)), reflect.ValueOf(func(v, exp uint16) bool {
		return v < exp
	}))

	RegisterLessThanOperator(reflect.TypeOf(int(0)), reflect.ValueOf(func(v int, exp int) bool {
		return v < exp
	}))

	RegisterLessThanOperator(reflect.TypeOf(int32(0)), reflect.ValueOf(func(v, exp int32) bool {
		return v < exp
	}))

	RegisterLessThanOperator(reflect.TypeOf(uint32(0)), reflect.ValueOf(func(v, exp uint32) bool {
		return v < exp
	}))

	RegisterLessThanOperator(reflect.TypeOf(int64(0)), reflect.ValueOf(func(v, exp int64) bool {
		return v < exp
	}))

	RegisterLessThanOperator(reflect.TypeOf(uint64(0)), reflect.ValueOf(func(v, exp uint64) bool {
		return v < exp
	}))

	RegisterGreaterThanOperator(reflect.TypeOf(byte(0)), reflect.ValueOf(func(v, exp byte) bool {
		return v > exp
	}))

	RegisterGreaterThanOperator(reflect.TypeOf(int8(0)), reflect.ValueOf(func(v, exp int8) bool {
		return v > exp
	}))

	RegisterGreaterThanOperator(reflect.TypeOf(uint8(0)), reflect.ValueOf(func(v, exp uint8) bool {
		return v > exp
	}))

	RegisterGreaterThanOperator(reflect.TypeOf(int16(0)), reflect.ValueOf(func(v, exp int16) bool {
		return v > exp
	}))

	RegisterGreaterThanOperator(reflect.TypeOf(uint16(0)), reflect.ValueOf(func(v, exp uint16) bool {
		return v > exp
	}))

	RegisterGreaterThanOperator(reflect.TypeOf(int(0)), reflect.ValueOf(func(v int, exp int) bool {
		return v > exp
	}))

	RegisterGreaterThanOperator(reflect.TypeOf(int32(0)), reflect.ValueOf(func(v, exp int32) bool {
		return v > exp
	}))

	RegisterGreaterThanOperator(reflect.TypeOf(uint32(0)), reflect.ValueOf(func(v, exp uint32) bool {
		return v > exp
	}))

	RegisterGreaterThanOperator(reflect.TypeOf(int64(0)), reflect.ValueOf(func(v, exp int64) bool {
		return v > exp
	}))

	RegisterGreaterThanOperator(reflect.TypeOf(uint64(0)), reflect.ValueOf(func(v, exp uint64) bool {
		return v > exp
	}))

	RegisterEqualsOperator(reflect.TypeOf(""), reflect.ValueOf(func(v, exp string) bool {
		return v == exp
	}))

	RegisterEqualsOperator(reflect.TypeOf([]byte(nil)), reflect.ValueOf(func(v, exp []byte) bool {
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

	RegisterEqualsOperator(reflect.TypeOf([]string(nil)), reflect.ValueOf(func(v, exp []string) bool {
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
