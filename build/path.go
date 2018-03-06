package build

import (
	"os"
)

type Path interface {
	Get() string
}

type PlainPath string

func (p PlainPath) Get() string {
	return string(p)
}

type EnvPath string

func (p EnvPath) Get() string {
	return os.ExpandEnv(string(p))
}

type LazyPath struct {
	*CachedTarget
}

func (p *LazyPath) Get() string {
	if p.CachedTarget.Output == nil {
		return ""
	}
	return p.CachedTarget.Output.Generated
}
