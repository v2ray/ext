package conf_test

import (
	"testing"

	. "v2ray.com/ext/assert"
	. "v2ray.com/ext/tools/conf"
)

func TestZeroBuffer(t *testing.T) {
	assert := With(t)

	bs := uint32(0)
	pConf := Policy{
		BufferSize: &bs,
	}
	p, err := pConf.Build()
	assert(err, IsNil)

	assert(p.Buffer.Enabled, IsFalse)
}
