package conf

import (
	"v2ray.com/core/app/policy"
)

type Timeout struct {
	Handshake      *uint32 `json:"handshake"`
	ConnectionIdle *uint32 `json:"connIdle"`
	UplinkOnly     *uint32 `json:"uplinkOnly"`
	DownlinkOnly   *uint32 `json:"downlinkOnly"`
}

func (t *Timeout) Build() (*policy.Policy_Timeout, error) {
	config := new(policy.Policy_Timeout)
	if t.Handshake != nil {
		config.Handshake = &policy.Second{Value: *t.Handshake}
	}
	if t.ConnectionIdle != nil {
		config.ConnectionIdle = &policy.Second{Value: *t.ConnectionIdle}
	}
	if t.UplinkOnly != nil {
		config.UplinkOnly = &policy.Second{Value: *t.UplinkOnly}
	}
	if t.DownlinkOnly != nil {
		config.DownlinkOnly = &policy.Second{Value: *t.DownlinkOnly}
	}
	return config, nil
}

type Policy struct {
	Timeout *Timeout `json:"timeout"`
}

func (p *Policy) Build() (*policy.Policy, error) {
	c := new(policy.Policy)
	if p.Timeout != nil {
		timeout, err := p.Timeout.Build()
		if err != nil {
			return nil, err
		}
		c.Timeout = timeout
	}
	return c, nil
}

type PolicyConfig struct {
	Levels map[uint32]*Policy `json:"levels"`
}

func (c *PolicyConfig) Build() (*policy.Config, error) {
	levels := make(map[uint32]*policy.Policy)
	for l, p := range c.Levels {
		if p != nil {
			pp, err := p.Build()
			if err != nil {
				return nil, err
			}
			levels[l] = pp
		}
	}
	return &policy.Config{
		Level: levels,
	}, nil
}
