package conf

import (
	"strings"

	"v2ray.com/core/app/dns"
	"v2ray.com/core/common/net"
)

/*
type NameServerConfig struct {
	Address *Address `json:"address"`
	Domains []string `json:"domains"`
}

func (c *NameServerConfig) Build() (*dns.NameServer, error) {
	if c.Address == nil {
		return nil, newError("NameServer address is not specified.")
	}

	return nil, nil
}
*/

// DnsConfig is a JSON serializable object for dns.Config.
type DnsConfig struct {
	Servers  []*Address          `json:"servers"`
	Hosts    map[string]*Address `json:"hosts"`
	ClientIP *Address            `json:"clientIp"`
}

// Build implements Buildable
func (c *DnsConfig) Build() (*dns.Config, error) {
	config := new(dns.Config)
	config.NameServers = make([]*net.Endpoint, len(c.Servers))

	if c.ClientIP != nil {
		if !c.ClientIP.Family().IsIPv4() && !c.ClientIP.Family().IsIPv6() {
			return nil, newError("not an IP address:", c.ClientIP.String())
		}
		config.ClientIp = []byte(c.ClientIP.IP())
	}

	for idx, server := range c.Servers {
		config.NameServers[idx] = &net.Endpoint{
			Network: net.Network_UDP,
			Address: server.Build(),
			Port:    53,
		}
	}

	if c.Hosts != nil {
		for domain, ip := range c.Hosts {
			if ip.Family() == net.AddressFamilyDomain {
				return nil, newError("domain is not expected in DNS hosts: ", ip.Domain())
			}

			mapping := &dns.Config_HostMapping{
				Ip: [][]byte{[]byte(ip.IP())},
			}
			if strings.HasPrefix(domain, "domain:") {
				mapping.Type = dns.DomainMatchingType_Subdomain
				mapping.Domain = domain[7:]
			} else {
				mapping.Type = dns.DomainMatchingType_Full
				mapping.Domain = domain
			}

			config.StaticHosts = append(config.StaticHosts, mapping)
		}
	}

	return config, nil
}
