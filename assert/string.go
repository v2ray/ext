package assert

import (
	"reflect"
	"strings"
)

var HasSubstring = CreateMatcher(reflect.TypeOf(""), reflect.ValueOf(func(a, b string) bool {
	return strings.Contains(a, b)
}), 2, "contains substring")
