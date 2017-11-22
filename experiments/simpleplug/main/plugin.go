package main

import (
	"context"

	"v2ray.com/core/common"
)

type Config struct {
	Name string
}

type Instance struct {
	config *Config
}

func (t *Instance) Name() string {
	return t.config.Name
}

func init() {
	common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return &Instance{
			config: config.(*Config),
		}, nil
	})
}

func main() {

}
