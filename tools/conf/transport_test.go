package conf_test

import (
	"encoding/json"
	"testing"

	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/headers/http"
	"v2ray.com/core/transport/internet/headers/noop"
	"v2ray.com/core/transport/internet/kcp"
	"v2ray.com/core/transport/internet/tcp"
	"v2ray.com/core/transport/internet/websocket"
	. "v2ray.com/ext/assert"
	. "v2ray.com/ext/tools/conf"
)

func TestTransportConfig(t *testing.T) {
	assert := With(t)

	rawJson := `{
    "tcpSettings": {
      "header": {
        "type": "http",
        "request": {
          "version": "1.1",
          "method": "GET",
          "path": "/b",
          "headers": {
            "a": "b",
            "c": "d"
          }
        },
        "response": {
          "version": "1.0",
          "status": "404",
          "reason": "Not Found"
        }
      }
    },
    "kcpSettings": {
      "mtu": 1200,
      "header": {
        "type": "none"
      }
    },
    "wsSettings": {
      "path": "/t"
    }
  }`

	var transportSettingsConf TransportConfig
	assert(json.Unmarshal([]byte(rawJson), &transportSettingsConf), IsNil)

	ts, err := transportSettingsConf.Build()
	assert(err, IsNil)

	assert(len(ts.TransportSettings), Equals, 3)
	var settingsCount uint32
	for _, settingsWithProtocol := range ts.TransportSettings {
		rawSettings, err := settingsWithProtocol.Settings.GetInstance()
		assert(err, IsNil)
		switch settings := rawSettings.(type) {
		case *tcp.Config:
			settingsCount++
			assert(settingsWithProtocol.Protocol == internet.TransportProtocol_TCP, IsTrue)
			rawHeader, err := settings.HeaderSettings.GetInstance()
			assert(err, IsNil)
			header := rawHeader.(*http.Config)
			assert(header.Request.GetVersionValue(), Equals, "1.1")
			assert(header.Request.Uri[0], Equals, "/b")
			assert(header.Request.Method.Value, Equals, "GET")
			var va, vc string
			for _, h := range header.Request.Header {
				switch h.Name {
				case "a":
					va = h.Value[0]
				case "c":
					vc = h.Value[0]
				default:
					t.Error("Unknown header ", h.String())
				}
			}
			assert(va, Equals, "b")
			assert(vc, Equals, "d")
			assert(header.Response.Version.Value, Equals, "1.0")
			assert(header.Response.Status.Code, Equals, "404")
			assert(header.Response.Status.Reason, Equals, "Not Found")
		case *kcp.Config:
			settingsCount++
			assert(settingsWithProtocol.Protocol, Equals, internet.TransportProtocol_MKCP)
			assert(settings.GetMTUValue(), Equals, uint32(1200))
			rawHeader, err := settings.HeaderConfig.GetInstance()
			assert(err, IsNil)
			header := rawHeader.(*noop.Config)
			assert(header, IsNotNil)
		case *websocket.Config:
			settingsCount++
			assert(settingsWithProtocol.Protocol, Equals, internet.TransportProtocol_WebSocket)
			assert(settings.Path, Equals, "/t")
		default:
			t.Error("Unknown type of settings.")
		}
	}
	assert(settingsCount, Equals, uint32(3))
}
