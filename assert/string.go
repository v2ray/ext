package assert

import (
	"strings"
)

var HasSubstring = CreateMatcher(func(a, b string) bool {
	return strings.Contains(a, b)
}, "contains substring")

var HasSuffix = CreateMatcher(func(a, b string) bool {
	return strings.HasSuffix(a, b)
}, "has suffix")

var HasPrefix = CreateMatcher(func(a, b string) bool {
	return strings.HasPrefix(a, b)
}, "has prefix")
