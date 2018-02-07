package conf

import (
	"strings"

	"v2ray.com/core/app/commander"
	"v2ray.com/core/app/proxyman/command"
	"v2ray.com/core/common/serial"
)

type ApiConfig struct {
	Tag      string   `json:"tag"`
	Services []string `json:"services"`
}

func (c *ApiConfig) Build() (*commander.Config, []*serial.TypedMessage, error) {
	if len(c.Tag) == 0 {
		return nil, nil, newError("Api tag can't be empty.")
	}

	services := make([]*serial.TypedMessage, 0, 16)
	for _, s := range c.Services {
		switch strings.ToLower(s) {
		case "handlerservice":
			services = append(services, serial.ToTypedMessage(&command.Config{}))
		}
	}

	return &commander.Config{
		Tag: c.Tag,
	}, services, nil
}
