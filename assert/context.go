package assert

import (
	"context"
)

var HasDone = CreateMatcher(func(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}, "has done.")
