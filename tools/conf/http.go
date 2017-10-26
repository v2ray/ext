package conf

import (
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/http"
)

type HttpAccount struct {
	Username string `json:"user"`
	Password string `json:"pass"`
}

type HttpServerConfig struct {
	Timeout  uint32         `json:"timeout"`
	Accounts []*HttpAccount `json:"accounts"`
}

func (c *HttpServerConfig) Build() (*serial.TypedMessage, error) {
	config := &http.ServerConfig{
		Timeout: c.Timeout,
	}

	if len(c.Accounts) > 0 {
		config.Accounts = make(map[string]string)
		for _, account := range c.Accounts {
			config.Accounts[account.Username] = account.Password
		}
	}

	return serial.ToTypedMessage(config), nil
}
