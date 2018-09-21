package conf_test

import (
	"encoding/json"
	"testing"

	"github.com/golang/protobuf/proto"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/transport"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/headers/http"
	"v2ray.com/core/transport/internet/headers/noop"
	"v2ray.com/core/transport/internet/kcp"
	"v2ray.com/core/transport/internet/tcp"
	"v2ray.com/core/transport/internet/websocket"
	. "v2ray.com/ext/tools/conf"
)

func TestTransportConfig(t *testing.T) {
	createParser := func() func(string) (proto.Message, error) {
		return func(s string) (proto.Message, error) {
			config := new(TransportConfig)
			if err := json.Unmarshal([]byte(s), config); err != nil {
				return nil, err
			}
			return config.Build()
		}
	}

	runMultiTestCase(t, []TestCase{
		{
			Input: `{
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
			}`,
			Parser: createParser(),
			Output: &transport.Config{
				TransportSettings: []*internet.TransportConfig{
					{
						Protocol: internet.TransportProtocol_TCP,
						Settings: serial.ToTypedMessage(&tcp.Config{
							HeaderSettings: serial.ToTypedMessage(&http.Config{
								Request: &http.RequestConfig{
									Version: &http.Version{Value: "1.1"},
									Method:  &http.Method{Value: "GET"},
									Uri:     []string{"/b"},
									Header: []*http.Header{
										{Name: "a", Value: []string{"b"}},
										{Name: "c", Value: []string{"d"}},
									},
								},
								Response: &http.ResponseConfig{
									Version: &http.Version{Value: "1.0"},
									Status:  &http.Status{Code: "404", Reason: "Not Found"},
									Header: []*http.Header{
										{
											Name:  "Content-Type",
											Value: []string{"application/octet-stream", "video/mpeg"},
										},
										{
											Name:  "Transfer-Encoding",
											Value: []string{"chunked"},
										},
										{
											Name:  "Connection",
											Value: []string{"keep-alive"},
										},
										{
											Name:  "Pragma",
											Value: []string{"no-cache"},
										},
										{
											Name:  "Cache-Control",
											Value: []string{"private", "no-cache"},
										},
									},
								},
							}),
						}),
					},
					{
						Protocol: internet.TransportProtocol_MKCP,
						Settings: serial.ToTypedMessage(&kcp.Config{
							Mtu:          &kcp.MTU{Value: 1200},
							HeaderConfig: serial.ToTypedMessage(&noop.Config{}),
						}),
					},
					{
						Protocol: internet.TransportProtocol_WebSocket,
						Settings: serial.ToTypedMessage(&websocket.Config{
							Path: "/t",
						}),
					},
				},
			},
		},
	})
}
