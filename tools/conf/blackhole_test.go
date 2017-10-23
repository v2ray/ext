package conf_test

import (
	"encoding/json"
	"testing"

	"v2ray.com/core/proxy/blackhole"
	. "v2ray.com/ext/assert"
	. "v2ray.com/ext/tools/conf"
)

func TestHTTPResponseJSON(t *testing.T) {
	assert := With(t)

	rawJson := `{
    "response": {
      "type": "http"
    }
  }`
	rawConfig := new(BlackholeConfig)
	err := json.Unmarshal([]byte(rawJson), rawConfig)
	assert(err, IsNil)

	ts, err := rawConfig.Build()
	assert(err, IsNil)
	iConfig, err := ts.GetInstance()
	assert(err, IsNil)
	config := iConfig.(*blackhole.Config)
	response, err := config.GetInternalResponse()
	assert(err, IsNil)

	_, ok := response.(*blackhole.HTTPResponse)
	assert(ok, IsTrue)
}
