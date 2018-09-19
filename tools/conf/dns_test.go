package conf_test

import (
	"encoding/json"
	"testing"

	"github.com/golang/protobuf/proto"
	"v2ray.com/core/app/dns"
	"v2ray.com/core/common/net"
	. "v2ray.com/ext/tools/conf"
)

func TestDnsConfigParsing(t *testing.T) {
	parserCreator := func() func(string) (proto.Message, error) {
		return func(s string) (proto.Message, error) {
			config := new(DnsConfig)
			if err := json.Unmarshal([]byte(s), config); err != nil {
				return nil, err
			}
			return config.Build()
		}
	}

	runMultiTestCase(t, []TestCase{
		{
			Input: `{
				"servers": [{
					"address": "8.8.8.8",
					"port": 5353,
					"domains": ["domain:v2ray.com"]
				}],
				"hosts": {
					"v2ray.com": "127.0.0.1"
				},
				"clientIp": "10.0.0.1"
			}`,
			Parser: parserCreator(),
			Output: &dns.Config{
				NameServer: []*dns.NameServer{
					{
						Address: &net.Endpoint{
							Address: &net.IPOrDomain{
								Address: &net.IPOrDomain_Ip{
									Ip: []byte{8, 8, 8, 8},
								},
							},
							Network: net.Network_UDP,
							Port:    5353,
						},
						PrioritizedDomain: []*dns.NameServer_PriorityDomain{
							{
								Type:   dns.DomainMatchingType_Subdomain,
								Domain: "v2ray.com",
							},
						},
					},
				},
				StaticHosts: []*dns.Config_HostMapping{
					{
						Type:   dns.DomainMatchingType_Full,
						Domain: "v2ray.com",
						Ip:     [][]byte{{127, 0, 0, 1}},
					},
				},
				ClientIp: []byte{10, 0, 0, 1},
			},
		},
	})
}
