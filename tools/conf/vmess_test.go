package conf_test

import (
	"encoding/json"
	"testing"

	"v2ray.com/core/common/protocol"
	"v2ray.com/core/proxy/vmess/outbound"
	. "v2ray.com/ext/assert"
	. "v2ray.com/ext/tools/conf"
)

func TestConfigTargetParsing(t *testing.T) {
	assert := With(t)

	rawJson := `{
    "vnext": [{
      "address": "127.0.0.1",
      "port": 80,
      "users": [
        {
          "id": "e641f5ad-9397-41e3-bf1a-e8740dfed019",
          "email": "love@v2ray.com",
          "level": 255
        }
      ]
    }]
  }`

	rawConfig := new(VMessOutboundConfig)
	err := json.Unmarshal([]byte(rawJson), &rawConfig)
	assert(err, IsNil)

	ts, err := rawConfig.Build()
	assert(err, IsNil)

	iConfig, err := ts.GetInstance()
	assert(err, IsNil)

	config := iConfig.(*outbound.Config)
	specPB := config.Receiver[0]
	spec := protocol.NewServerSpecFromPB(*specPB)
	assert(spec.Destination().String(), Equals, "tcp:127.0.0.1:80")
}
