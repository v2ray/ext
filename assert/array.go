package assert

var HasStringElement = CreateMatcher(func(a []string, b string) bool {
	for _, v := range a {
		if v == b {
			return true
		}
	}
	return false
}, "contains string element")
