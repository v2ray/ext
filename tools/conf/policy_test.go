package conf_test

import (
	"testing"

	. "v2ray.com/ext/assert"
	. "v2ray.com/ext/tools/conf"
)

func TestBufferSize(t *testing.T) {
	assert := With(t)

	cases := []struct {
		Input  int32
		Output int32
	}{
		{
			Input:  0,
			Output: 0,
		},
		{
			Input:  -1,
			Output: -1,
		},
		{
			Input:  1,
			Output: 1024,
		},
	}

	for _, c := range cases {
		bs := int32(c.Input)
		pConf := Policy{
			BufferSize: &bs,
		}
		p, err := pConf.Build()
		assert(err, IsNil)

		assert(p.Buffer.Connection, Equals, int32(c.Output))
	}

}
