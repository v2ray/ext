package conf

import (
	"v2ray.com/core/app/dns"
	v2net "v2ray.com/core/common/net"
)

// DnsConfig is a JSON serializable object for dns.Config.
type DnsConfig struct {
	Servers  []*Address          `json:"servers"`
	Hosts    map[string]*Address `json:"hosts"`
	ClientV4 *Address            `json:"clientIp"`
	ClientV6 *Address            `json:"clientIp6"`
}

// Build implements Buildable
func (c *DnsConfig) Build() (*dns.Config, error) {
	config := new(dns.Config)
	config.NameServers = make([]*v2net.Endpoint, len(c.Servers))

	if c.ClientV4 != nil {
		if c.ClientV4.Family() != v2net.AddressFamilyIPv4 {
			return nil, newError("not an IPV4 address:", c.ClientV4.String())
		}
		if config.ClientIp == nil {
			config.ClientIp = &dns.Config_ClientIP{}
		}
		config.ClientIp.V4 = []byte(c.ClientV4.IP())
	}

	if c.ClientV6 != nil {
		if c.ClientV6.Family() != v2net.AddressFamilyIPv6 {
			return nil, newError("not an IPV6 address:", c.ClientV4.String())
		}
		if config.ClientIp == nil {
			config.ClientIp = &dns.Config_ClientIP{}
		}
		config.ClientIp.V6 = []byte(c.ClientV6.IP())
	}

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

	return config, nil
}
