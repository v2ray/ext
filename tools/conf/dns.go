package conf

import (
	"v2ray.com/core/app/dns"
	v2net "v2ray.com/core/common/net"
)

// DnsConfig is a JSON serializable object for dns.Config.
type DnsConfig struct {
	Servers []*Address          `json:"servers"`
	Hosts   map[string]*Address `json:"hosts"`
}

// Build implements Buildable
func (c *DnsConfig) Build() *dns.Config {
	config := new(dns.Config)
	config.NameServers = make([]*v2net.Endpoint, len(c.Servers))
	for idx, server := range c.Servers {
		config.NameServers[idx] = &v2net.Endpoint{
			Network: v2net.Network_UDP,
			Address: server.Build(),
			Port:    53,
		}
	}

	if c.Hosts != nil {
		config.Hosts = make(map[string]*v2net.IPOrDomain)
		for domain, ip := range c.Hosts {
			config.Hosts[domain] = ip.Build()
		}
	}

	return config
}
