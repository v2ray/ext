package conf_test

import (
	"encoding/json"
	"testing"

	"v2ray.com/core/proxy/shadowsocks"
	. "v2ray.com/ext/assert"
	. "v2ray.com/ext/tools/conf"
)

func TestShadowsocksServerConfigParsing(t *testing.T) {
	assert := With(t)

	rawJson := `{
    "method": "aes-128-cfb",
    "password": "v2ray-password"
  }`

	rawConfig := new(ShadowsocksServerConfig)
	err := json.Unmarshal([]byte(rawJson), rawConfig)
	assert(err, IsNil)

	ts, err := rawConfig.Build()
	assert(err, IsNil)
	iConfig, err := ts.GetInstance()
	assert(err, IsNil)
	config := iConfig.(*shadowsocks.ServerConfig)

	rawAccount, err := config.User.GetTypedAccount()
	assert(err, IsNil)

	account, ok := rawAccount.(*shadowsocks.MemoryAccount)
	assert(ok, IsTrue)

	assert(account.Cipher.KeySize(), Equals, int32(16))
	assert(account.Key, Equals, []byte{160, 224, 26, 2, 22, 110, 9, 80, 65, 52, 80, 20, 38, 243, 224, 241})
}
