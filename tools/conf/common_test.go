package conf_test

import (
	"encoding/json"
	"os"
	"testing"

	"v2ray.com/core/common/net"
	. "v2ray.com/ext/assert"
	. "v2ray.com/ext/tools/conf"
)

func TestStringListUnmarshalError(t *testing.T) {
	assert := With(t)

	rawJson := `1234`
	list := new(StringList)
	err := json.Unmarshal([]byte(rawJson), list)
	assert(err, IsNotNil)
}

func TestStringListLen(t *testing.T) {
	assert := With(t)

	rawJson := `"a, b, c, d"`
	list := new(StringList)
	err := json.Unmarshal([]byte(rawJson), list)
	assert(err, IsNil)
	assert(list.Len(), Equals, 4)
}

func TestIPParsing(t *testing.T) {
	assert := With(t)

	rawJson := "\"8.8.8.8\""
	var address Address
	err := json.Unmarshal([]byte(rawJson), &address)
	assert(err, IsNil)
	assert([]byte(address.IP()), Equals, []byte{8, 8, 8, 8})
}

func TestDomainParsing(t *testing.T) {
	assert := With(t)

	rawJson := "\"v2ray.com\""
	var address Address
	err := json.Unmarshal([]byte(rawJson), &address)
	assert(err, IsNil)
	assert(address.Domain(), Equals, "v2ray.com")
}

func TestInvalidAddressJson(t *testing.T) {
	assert := With(t)

	rawJson := "1234"
	var address Address
	err := json.Unmarshal([]byte(rawJson), &address)
	assert(err, IsNotNil)
}

func TestStringNetwork(t *testing.T) {
	assert := With(t)

	var network Network
	err := json.Unmarshal([]byte(`"tcp"`), &network)
	assert(err, IsNil)
	assert(network.Build() == net.Network_TCP, IsTrue)
}

func TestArrayNetworkList(t *testing.T) {
	assert := With(t)

	var list NetworkList
	err := json.Unmarshal([]byte("[\"Tcp\"]"), &list)
	assert(err, IsNil)

	nlist := list.Build()
	assert(net.HasNetwork(nlist, net.Network_TCP), IsTrue)
	assert(net.HasNetwork(nlist, net.Network_UDP), IsFalse)
}

func TestStringNetworkList(t *testing.T) {
	assert := With(t)

	var list NetworkList
	err := json.Unmarshal([]byte("\"TCP, ip\""), &list)
	assert(err, IsNil)

	nlist := list.Build()
	assert(net.HasNetwork(nlist, net.Network_TCP), IsTrue)
	assert(net.HasNetwork(nlist, net.Network_UDP), IsFalse)
}

func TestInvalidNetworkJson(t *testing.T) {
	assert := With(t)

	var list NetworkList
	err := json.Unmarshal([]byte("0"), &list)
	assert(err, IsNotNil)
}

func TestIntPort(t *testing.T) {
	assert := With(t)

	var portRange PortRange
	err := json.Unmarshal([]byte("1234"), &portRange)
	assert(err, IsNil)

	assert(portRange.From, Equals, uint32(1234))
	assert(portRange.To, Equals, uint32(1234))
}

func TestOverRangeIntPort(t *testing.T) {
	assert := With(t)

	var portRange PortRange
	err := json.Unmarshal([]byte("70000"), &portRange)
	assert(err, IsNotNil)

	err = json.Unmarshal([]byte("-1"), &portRange)
	assert(err, IsNotNil)
}

func TestEnvPort(t *testing.T) {
	assert := With(t)

	assert(os.Setenv("PORT", "1234"), IsNil)

	var portRange PortRange
	err := json.Unmarshal([]byte("\"env:PORT\""), &portRange)
	assert(err, IsNil)

	assert(portRange.From, Equals, uint32(1234))
	assert(portRange.To, Equals, uint32(1234))
}

func TestSingleStringPort(t *testing.T) {
	assert := With(t)

	var portRange PortRange
	err := json.Unmarshal([]byte("\"1234\""), &portRange)
	assert(err, IsNil)

	assert(portRange.From, Equals, uint32(1234))
	assert(portRange.To, Equals, uint32(1234))
}

func TestStringPairPort(t *testing.T) {
	assert := With(t)

	var portRange PortRange
	err := json.Unmarshal([]byte("\"1234-5678\""), &portRange)
	assert(err, IsNil)

	assert(portRange.From, Equals, uint32(1234))
	assert(portRange.To, Equals, uint32(5678))
}

func TestOverRangeStringPort(t *testing.T) {
	assert := With(t)

	var portRange PortRange
	err := json.Unmarshal([]byte("\"65536\""), &portRange)
	assert(err, IsNotNil)

	err = json.Unmarshal([]byte("\"70000-80000\""), &portRange)
	assert(err, IsNotNil)

	err = json.Unmarshal([]byte("\"1-90000\""), &portRange)
	assert(err, IsNotNil)

	err = json.Unmarshal([]byte("\"700-600\""), &portRange)
	assert(err, IsNotNil)
}

func TestUserParsing(t *testing.T) {
	assert := With(t)

	user := new(User)
	err := json.Unmarshal([]byte(`{
    "id": "96edb838-6d68-42ef-a933-25f7ac3a9d09",
    "email": "love@v2ray.com",
    "level": 1,
    "alterId": 100
  }`), user)
	assert(err, IsNil)

	nUser := user.Build()
	assert(byte(nUser.Level), Equals, byte(1))
	assert(nUser.Email, Equals, "love@v2ray.com")
}

func TestInvalidUserJson(t *testing.T) {
	assert := With(t)

	user := new(User)
	err := json.Unmarshal([]byte(`{"email": 1234}`), user)
	assert(err, IsNotNil)
}
