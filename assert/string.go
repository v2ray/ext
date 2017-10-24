package assert

import (
	"strings"
)

var HasSubstring = CreateMatcher(func(a, b string) bool {
	return strings.Contains(a, b)
}, "contains substring")
