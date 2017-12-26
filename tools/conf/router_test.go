package conf_test

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"v2ray.com/ext/sysio"

	"v2ray.com/core/common/net"
	"v2ray.com/core/common/platform"
	"v2ray.com/core/proxy"
	. "v2ray.com/ext/assert"
	. "v2ray.com/ext/tools/conf"
)

func makeDestination(ip string) net.Destination {
	return net.TCPDestination(net.IPAddress(net.ParseIP(ip)), 80)
}

func makeDomainDestination(domain string) net.Destination {
	return net.TCPDestination(net.DomainAddress(domain), 80)
}

func TestChinaIPJson(t *testing.T) {
	assert := With(t)

	fileBytes, err := sysio.ReadFile(filepath.Join(os.Getenv("GOPATH"), "src", "v2ray.com", "core", "release", "config", "geoip.dat"))
	assert(err, IsNil)

	assert(ioutil.WriteFile(platform.GetAssetLocation("geoip.dat"), fileBytes, 0666), IsNil)

	rule, err := ParseRule([]byte(`{
    "type": "chinaip",
    "outboundTag": "x"
	}`))
	assert(err, IsNil)
	assert(rule.Tag, Equals, "x")
	cond, err := rule.BuildCondition()
	assert(err, IsNil)
	assert(cond.Apply(proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.ParseAddress("121.14.1.189"), 80))), IsTrue)    // sina.com.cn
	assert(cond.Apply(proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.ParseAddress("101.226.103.106"), 80))), IsTrue) // qq.com
	assert(cond.Apply(proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.ParseAddress("115.239.210.36"), 80))), IsTrue)  // image.baidu.com
	assert(cond.Apply(proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.ParseAddress("120.135.126.1"), 80))), IsTrue)

	assert(cond.Apply(proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.ParseAddress("8.8.8.8"), 80))), IsFalse)
}

func TestChinaSitesJson(t *testing.T) {
	assert := With(t)

	fileBytes, err := sysio.ReadFile(filepath.Join(os.Getenv("GOPATH"), "src", "v2ray.com", "core", "release", "config", "geosite.dat"))
	assert(err, IsNil)

	assert(ioutil.WriteFile(platform.GetAssetLocation("geosite.dat"), fileBytes, 0666), IsNil)

	rule, err := ParseRule([]byte(`{
    "type": "chinasites",
    "outboundTag": "y"
  }`))
	assert(err, IsNil)
	assert(rule.Tag, Equals, "y")
	cond, err := rule.BuildCondition()
	assert(err, IsNil)
	assert(cond.Apply(proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.ParseAddress("v.qq.com"), 80))), IsTrue)
	assert(cond.Apply(proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.ParseAddress("www.163.com"), 80))), IsTrue)
	assert(cond.Apply(proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.ParseAddress("ngacn.cc"), 80))), IsTrue)
	assert(cond.Apply(proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.ParseAddress("12306.cn"), 80))), IsTrue)

	assert(cond.Apply(proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.ParseAddress("v2ray.com"), 80))), IsFalse)
}

func TestDomainRule(t *testing.T) {
	assert := With(t)

	rule, err := ParseRule([]byte(`{
    "type": "field",
    "domain": [
      "ooxx.com",
      "oxox.com",
      "regexp:\\.cn$"
    ],
    "network": "tcp",
    "outboundTag": "direct"
  }`))
	assert(err, IsNil)
	assert(rule, IsNotNil)
	cond, err := rule.BuildCondition()
	assert(err, IsNil)
	assert(cond.Apply(proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.ParseAddress("www.ooxx.com"), 80))), IsTrue)
	assert(cond.Apply(proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.ParseAddress("www.aabb.com"), 80))), IsFalse)
	assert(cond.Apply(proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.IPAddress([]byte{127, 0, 0, 1}), 80))), IsFalse)
	assert(cond.Apply(proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.ParseAddress("www.12306.cn"), 80))), IsTrue)
	assert(cond.Apply(proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.ParseAddress("www.acn.com"), 80))), IsFalse)
}

func TestIPRule(t *testing.T) {
	assert := With(t)

	rule, err := ParseRule([]byte(`{
    "type": "field",
    "ip": [
      "10.0.0.0/8",
      "192.0.0.0/24"
    ],
    "network": "tcp",
    "outboundTag": "direct"
  }`))
	assert(err, IsNil)
	assert(rule, IsNotNil)
	cond, err := rule.BuildCondition()
	assert(err, IsNil)
	assert(cond.Apply(proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.DomainAddress("www.ooxx.com"), 80))), IsFalse)
	assert(cond.Apply(proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.IPAddress([]byte{10, 0, 0, 1}), 80))), IsTrue)
	assert(cond.Apply(proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.IPAddress([]byte{127, 0, 0, 1}), 80))), IsFalse)
	assert(cond.Apply(proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.IPAddress([]byte{192, 0, 0, 1}), 80))), IsTrue)
}

func TestSourceIPRule(t *testing.T) {
	assert := With(t)

	rule, err := ParseRule([]byte(`{
    "type": "field",
    "source": [
      "10.0.0.0/8",
      "192.0.0.0/24"
    ],
    "outboundTag": "direct"
  }`))
	assert(err, IsNil)
	assert(rule, IsNotNil)
	cond, err := rule.BuildCondition()
	assert(err, IsNil)
	assert(cond.Apply(proxy.ContextWithSource(context.Background(), net.TCPDestination(net.DomainAddress("www.ooxx.com"), 80))), IsFalse)
	assert(cond.Apply(proxy.ContextWithSource(context.Background(), net.TCPDestination(net.IPAddress([]byte{10, 0, 0, 1}), 80))), IsTrue)
	assert(cond.Apply(proxy.ContextWithSource(context.Background(), net.TCPDestination(net.IPAddress([]byte{127, 0, 0, 1}), 80))), IsFalse)
	assert(cond.Apply(proxy.ContextWithSource(context.Background(), net.TCPDestination(net.IPAddress([]byte{192, 0, 0, 1}), 80))), IsTrue)
}
