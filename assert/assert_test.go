package assert_test

import (
	"testing"

	. "v2ray.com/ext/assert"
)

func TestIntEquals(t *testing.T) {
	assert := With(t)

	assert(1, Equals, 1)
}

func TestStringEquals(t *testing.T) {
	assert := With(t)

	assert("abcd", Equals, "abcd")
}

func TestByteArrayNotEquals(t *testing.T) {
	assert := With(t)

	assert([]byte{1, 2, 3, 4}, NotEquals, []byte{1, 2, 3})
	assert([]byte{1, 2, 3, 4}, NotEquals, []byte{1, 2, 3, 5})
}

func TestNil(t *testing.T) {
	assert := With(t)

	var err error
	assert(err, IsNil)
}

func TestPanic(t *testing.T) {
	assert := With(t)

	assert(func() { panic("panic on purpose.") }, Panics)
}
