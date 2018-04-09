package conf

import (
	"v2ray.com/core/transport"
	"v2ray.com/core/transport/internet"
)

type TransportConfig struct {
	TCPConfig  *TCPConfig          `json:"tcpSettings"`
	KCPConfig  *KCPConfig          `json:"kcpSettings"`
	WSConfig   *WebSocketConfig    `json:"wsSettings"`
	HTTPConfig *HTTPConfig         `json:"httpSettings"`
	DSConfig   *DomainSocketConfig `json:"dsSettings"`
}

// Build implements Buildable.
func (c *TransportConfig) Build() (*transport.Config, error) {
	config := new(transport.Config)

	if c.TCPConfig != nil {
		ts, err := c.TCPConfig.Build()
		if err != nil {
			return nil, newError("failed to build TCP config").Base(err).AtError()
		}
		config.TransportSettings = append(config.TransportSettings, &internet.TransportConfig{
			Protocol: internet.TransportProtocol_TCP,
			Settings: ts,
		})
	}

	if c.KCPConfig != nil {
		ts, err := c.KCPConfig.Build()
		if err != nil {
			return nil, newError("failed to build mKCP config").Base(err).AtError()
		}
		config.TransportSettings = append(config.TransportSettings, &internet.TransportConfig{
			Protocol: internet.TransportProtocol_MKCP,
			Settings: ts,
		})
	}

	if c.WSConfig != nil {
		ts, err := c.WSConfig.Build()
		if err != nil {
			return nil, newError("failed to build WebSocket config").Base(err)
		}
		config.TransportSettings = append(config.TransportSettings, &internet.TransportConfig{
			Protocol: internet.TransportProtocol_WebSocket,
			Settings: ts,
		})
	}

	if c.HTTPConfig != nil {
		ts, err := c.HTTPConfig.Build()
		if err != nil {
			return nil, newError("Failed to build HTTP config.").Base(err)
		}
		config.TransportSettings = append(config.TransportSettings, &internet.TransportConfig{
			Protocol: internet.TransportProtocol_HTTP,
			Settings: ts,
		})
	}

	if c.DSConfig != nil {
		ds, err := c.DSConfig.Build()
		if err != nil {
			return nil, newError("Failed to build DomainSocket config.").Base(err)
		}
		config.TransportSettings = append(config.TransportSettings, &internet.TransportConfig{
			Protocol: internet.TransportProtocol_DomainSocket,
			Settings: ds,
		})
	}

	return config, nil
}
