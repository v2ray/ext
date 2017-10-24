package conf_test

import (
	"encoding/json"
	"testing"

	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/proxy/socks"
	. "v2ray.com/ext/assert"
	. "v2ray.com/ext/tools/conf"
)

func TestSocksInboundConfig(t *testing.T) {
	assert := With(t)

	rawJson := `{
    "auth": "password",
    "accounts": [
      {
        "user": "my-username",
        "pass": "my-password"
      }
    ],
    "udp": false,
    "ip": "127.0.0.1",
    "timeout": 5
  }`

	config := new(SocksServerConfig)
	err := json.Unmarshal([]byte(rawJson), &config)
	assert(err, IsNil)

	message, err := config.Build()
	assert(err, IsNil)

	iConfig, err := message.GetInstance()
	assert(err, IsNil)

	socksConfig := iConfig.(*socks.ServerConfig)
	assert(socksConfig.AuthType == socks.AuthType_PASSWORD, IsTrue)
	assert(len(socksConfig.Accounts), Equals, 1)
	assert(socksConfig.Accounts["my-username"], Equals, "my-password")
	assert(socksConfig.UdpEnabled, IsFalse)
	assert(socksConfig.Address.AsAddress().String(), Equals, net.LocalHostIP.String())
	assert(socksConfig.Timeout, Equals, uint32(5))
}

func TestSocksOutboundConfig(t *testing.T) {
	assert := With(t)

	rawJson := `{
    "servers": [{
      "address": "127.0.0.1",
      "port": 1234,
      "users": [
        {"user": "test user", "pass": "test pass", "email": "test@email.com"}
      ]
    }]
  }`

	config := new(SocksClientConfig)
	err := json.Unmarshal([]byte(rawJson), &config)
	assert(err, IsNil)

	message, err := config.Build()
	assert(err, IsNil)

	iConfig, err := message.GetInstance()
	assert(err, IsNil)

	socksConfig := iConfig.(*socks.ClientConfig)
	assert(len(socksConfig.Server), Equals, 1)

	ss := protocol.NewServerSpecFromPB(*socksConfig.Server[0])
	assert(ss.Destination().String(), Equals, "tcp:127.0.0.1:1234")

	user := ss.PickUser()
	assert(user.Email, Equals, "test@email.com")

	account, err := user.GetTypedAccount()
	assert(err, IsNil)

	socksAccount := account.(*socks.Account)
	assert(socksAccount.Username, Equals, "test user")
	assert(socksAccount.Password, Equals, "test pass")
}
