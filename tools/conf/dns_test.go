package conf_test

import (
	"encoding/json"
	"testing"

	"v2ray.com/core/common/net"
	. "v2ray.com/core/common/net/testing"
	. "v2ray.com/ext/assert"
	. "v2ray.com/ext/tools/conf"
)

func TestDnsConfigParsing(t *testing.T) {
	assert := With(t)

	rawJson := `{
    "servers": ["8.8.8.8"]
  }`

	jsonConfig := new(DnsConfig)
	err := json.Unmarshal([]byte(rawJson), jsonConfig)
	assert(err, IsNil)

	config, err := jsonConfig.Build()
	assert(err, IsNil)
	assert(len(config.NameServers), Equals, 1)
	dest := config.NameServers[0].AsDestination()
	assert(dest, IsUDP)
	assert(dest.Address.String(), Equals, "8.8.8.8")
	assert(dest.Port, Equals, net.Port(53))
}
