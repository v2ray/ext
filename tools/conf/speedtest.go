package conf

import (
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/speedtest"
)

type SpeedTestConfig struct{}

func (c *SpeedTestConfig) Build() (*serial.TypedMessage, error) {
	config := new(speedtest.Config)
	return serial.ToTypedMessage(config), nil
}
