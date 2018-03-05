package build_test

import (
	"testing"

	. "v2ray.com/ext/assert"
	. "v2ray.com/ext/build"
)

func TestBuildEnvs(t *testing.T) {
	assert := With(t)

	target := &GoTarget{
		OS:   Windows,
		Arch: Amd64,
	}

	envs := target.Envs()
	assert(envs, HasStringElement, "GOOS=windows")
	assert(envs, HasStringElement, "GOARCH=amd64")
}
