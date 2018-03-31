package conf

import (
	"encoding/json"
	"strings"

	"v2ray.com/core"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/app/stats"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
)

var (
	inboundConfigLoader = NewJSONConfigLoader(ConfigCreatorCache{
		"dokodemo-door": func() interface{} { return new(DokodemoConfig) },
		"http":          func() interface{} { return new(HttpServerConfig) },
		"shadowsocks":   func() interface{} { return new(ShadowsocksServerConfig) },
		"socks":         func() interface{} { return new(SocksServerConfig) },
		"vmess":         func() interface{} { return new(VMessInboundConfig) },
	}, "protocol", "settings")

	outboundConfigLoader = NewJSONConfigLoader(ConfigCreatorCache{
		"blackhole":   func() interface{} { return new(BlackholeConfig) },
		"freedom":     func() interface{} { return new(FreedomConfig) },
		"shadowsocks": func() interface{} { return new(ShadowsocksClientConfig) },
		"vmess":       func() interface{} { return new(VMessOutboundConfig) },
		"socks":       func() interface{} { return new(SocksClientConfig) },
	}, "protocol", "settings")
)

func toProtocolList(s []string) ([]proxyman.KnownProtocols, error) {
	kp := make([]proxyman.KnownProtocols, 0, 8)
	for _, p := range s {
		switch strings.ToLower(p) {
		case "http":
			kp = append(kp, proxyman.KnownProtocols_HTTP)
		case "https", "tls", "ssl":
			kp = append(kp, proxyman.KnownProtocols_TLS)
		default:
			return nil, newError("Unknown protocol: ", p)
		}
	}
	return kp, nil
}

type InboundConnectionConfig struct {
	Port           uint16          `json:"port"`
	Listen         *Address        `json:"listen"`
	Protocol       string          `json:"protocol"`
	StreamSetting  *StreamConfig   `json:"streamSettings"`
	Settings       json.RawMessage `json:"settings"`
	Tag            string          `json:"tag"`
	DomainOverride *StringList     `json:"domainOverride"`
}

// Build implements Buildable.
func (c *InboundConnectionConfig) Build() (*core.InboundHandlerConfig, error) {
	receiverConfig := &proxyman.ReceiverConfig{
		PortRange: &v2net.PortRange{
			From: uint32(c.Port),
			To:   uint32(c.Port),
		},
	}
	if c.Listen != nil {
		if c.Listen.Family().IsDomain() {
			return nil, newError("unable to listen on domain address: " + c.Listen.Domain())
		}
		receiverConfig.Listen = c.Listen.Build()
	}
	if c.StreamSetting != nil {
		ts, err := c.StreamSetting.Build()
		if err != nil {
			return nil, err
		}
		receiverConfig.StreamSettings = ts
	}
	if c.DomainOverride != nil {
		kp, err := toProtocolList(*c.DomainOverride)
		if err != nil {
			return nil, newError("failed to parse inbound config").Base(err)
		}
		receiverConfig.DomainOverride = kp
	}

	jsonConfig, err := inboundConfigLoader.LoadWithID(c.Settings, c.Protocol)
	if err != nil {
		return nil, newError("failed to load inbound config.").Base(err)
	}
	if dokodemoConfig, ok := jsonConfig.(*DokodemoConfig); ok {
		receiverConfig.ReceiveOriginalDestination = dokodemoConfig.Redirect
	}
	ts, err := jsonConfig.(Buildable).Build()
	if err != nil {
		return nil, err
	}

	return &core.InboundHandlerConfig{
		Tag:              c.Tag,
		ReceiverSettings: serial.ToTypedMessage(receiverConfig),
		ProxySettings:    ts,
	}, nil
}

type MuxConfig struct {
	Enabled     bool   `json:"enabled"`
	Concurrency uint16 `json:"concurrency"`
}

func (c *MuxConfig) GetConcurrency() uint16 {
	if c.Concurrency == 0 {
		return 8
	}
	return c.Concurrency
}

type OutboundConnectionConfig struct {
	Protocol      string          `json:"protocol"`
	SendThrough   *Address        `json:"sendThrough"`
	StreamSetting *StreamConfig   `json:"streamSettings"`
	ProxySettings *ProxyConfig    `json:"proxySettings"`
	Settings      json.RawMessage `json:"settings"`
	Tag           string          `json:"tag"`
	MuxSettings   *MuxConfig      `json:"mux"`
}

// Build implements Buildable.
func (c *OutboundConnectionConfig) Build() (*core.OutboundHandlerConfig, error) {
	senderSettings := &proxyman.SenderConfig{}

	if c.SendThrough != nil {
		address := c.SendThrough
		if address.Family().IsDomain() {
			return nil, newError("invalid sendThrough address: " + address.String())
		}
		senderSettings.Via = address.Build()
	}
	if c.StreamSetting != nil {
		ss, err := c.StreamSetting.Build()
		if err != nil {
			return nil, err
		}
		senderSettings.StreamSettings = ss
	}
	if c.ProxySettings != nil {
		ps, err := c.ProxySettings.Build()
		if err != nil {
			return nil, newError("invalid outbound proxy settings").Base(err)
		}
		senderSettings.ProxySettings = ps
	}

	if c.MuxSettings != nil && c.MuxSettings.Enabled {
		senderSettings.MultiplexSettings = &proxyman.MultiplexingConfig{
			Enabled:     true,
			Concurrency: uint32(c.MuxSettings.GetConcurrency()),
		}
	}

	rawConfig, err := outboundConfigLoader.LoadWithID(c.Settings, c.Protocol)
	if err != nil {
		return nil, newError("failed to parse outbound config").Base(err)
	}
	ts, err := rawConfig.(Buildable).Build()
	if err != nil {
		return nil, err
	}

	return &core.OutboundHandlerConfig{
		SenderSettings: serial.ToTypedMessage(senderSettings),
		ProxySettings:  ts,
		Tag:            c.Tag,
	}, nil
}

type InboundDetourAllocationConfig struct {
	Strategy    string  `json:"strategy"`
	Concurrency *uint32 `json:"concurrency"`
	RefreshMin  *uint32 `json:"refresh"`
}

// Build implements Buildable.
func (c *InboundDetourAllocationConfig) Build() (*proxyman.AllocationStrategy, error) {
	config := new(proxyman.AllocationStrategy)
	switch strings.ToLower(c.Strategy) {
	case "always":
		config.Type = proxyman.AllocationStrategy_Always
	case "random":
		config.Type = proxyman.AllocationStrategy_Random
	case "external":
		config.Type = proxyman.AllocationStrategy_External
	default:
		return nil, newError("unknown allocation strategy: ", c.Strategy)
	}
	if c.Concurrency != nil {
		config.Concurrency = &proxyman.AllocationStrategy_AllocationStrategyConcurrency{
			Value: *c.Concurrency,
		}
	}

	if c.RefreshMin != nil {
		config.Refresh = &proxyman.AllocationStrategy_AllocationStrategyRefresh{
			Value: *c.RefreshMin,
		}
	}

	return config, nil
}

type InboundDetourConfig struct {
	Protocol       string                         `json:"protocol"`
	PortRange      *PortRange                     `json:"port"`
	ListenOn       *Address                       `json:"listen"`
	Settings       json.RawMessage                `json:"settings"`
	Tag            string                         `json:"tag"`
	Allocation     *InboundDetourAllocationConfig `json:"allocate"`
	StreamSetting  *StreamConfig                  `json:"streamSettings"`
	DomainOverride *StringList                    `json:"domainOverride"`
}

// Build implements Buildable.
func (c *InboundDetourConfig) Build() (*core.InboundHandlerConfig, error) {
	receiverSettings := &proxyman.ReceiverConfig{}

	if c.PortRange == nil {
		return nil, newError("port range not specified in InboundDetour.")
	}
	receiverSettings.PortRange = c.PortRange.Build()

	if c.ListenOn != nil {
		if c.ListenOn.Family().IsDomain() {
			return nil, newError("unable to listen on domain address: ", c.ListenOn.Domain())
		}
		receiverSettings.Listen = c.ListenOn.Build()
	}
	if c.Allocation != nil {
		as, err := c.Allocation.Build()
		if err != nil {
			return nil, err
		}
		receiverSettings.AllocationStrategy = as
	}
	if c.StreamSetting != nil {
		ss, err := c.StreamSetting.Build()
		if err != nil {
			return nil, err
		}
		receiverSettings.StreamSettings = ss
	}
	if c.DomainOverride != nil {
		kp, err := toProtocolList(*c.DomainOverride)
		if err != nil {
			return nil, newError("failed to parse inbound detour config").Base(err)
		}
		receiverSettings.DomainOverride = kp
	}

	rawConfig, err := inboundConfigLoader.LoadWithID(c.Settings, c.Protocol)
	if err != nil {
		return nil, newError("failed to load inbound detour config.").Base(err)
	}
	if dokodemoConfig, ok := rawConfig.(*DokodemoConfig); ok {
		receiverSettings.ReceiveOriginalDestination = dokodemoConfig.Redirect
	}
	ts, err := rawConfig.(Buildable).Build()
	if err != nil {
		return nil, err
	}

	return &core.InboundHandlerConfig{
		Tag:              c.Tag,
		ReceiverSettings: serial.ToTypedMessage(receiverSettings),
		ProxySettings:    ts,
	}, nil
}

type OutboundDetourConfig struct {
	Protocol      string          `json:"protocol"`
	SendThrough   *Address        `json:"sendThrough"`
	Tag           string          `json:"tag"`
	Settings      json.RawMessage `json:"settings"`
	StreamSetting *StreamConfig   `json:"streamSettings"`
	ProxySettings *ProxyConfig    `json:"proxySettings"`
	MuxSettings   *MuxConfig      `json:"mux"`
}

// Build implements Buildable.
func (c *OutboundDetourConfig) Build() (*core.OutboundHandlerConfig, error) {
	senderSettings := &proxyman.SenderConfig{}

	if c.SendThrough != nil {
		address := c.SendThrough
		if address.Family().IsDomain() {
			return nil, newError("unable to send through: " + address.String())
		}
		senderSettings.Via = address.Build()
	}

	if c.StreamSetting != nil {
		ss, err := c.StreamSetting.Build()
		if err != nil {
			return nil, err
		}
		senderSettings.StreamSettings = ss
	}

	if c.ProxySettings != nil {
		ps, err := c.ProxySettings.Build()
		if err != nil {
			return nil, newError("invalid outbound detour proxy settings.").Base(err)
		}
		senderSettings.ProxySettings = ps
	}

	if c.MuxSettings != nil && c.MuxSettings.Enabled {
		senderSettings.MultiplexSettings = &proxyman.MultiplexingConfig{
			Enabled:     true,
			Concurrency: uint32(c.MuxSettings.GetConcurrency()),
		}
	}

	rawConfig, err := outboundConfigLoader.LoadWithID(c.Settings, c.Protocol)
	if err != nil {
		return nil, newError("failed to parse to outbound detour config.").Base(err)
	}
	ts, err := rawConfig.(Buildable).Build()
	if err != nil {
		return nil, err
	}

	return &core.OutboundHandlerConfig{
		SenderSettings: serial.ToTypedMessage(senderSettings),
		Tag:            c.Tag,
		ProxySettings:  ts,
	}, nil
}

type StatsConfig struct{}

func (c *StatsConfig) Build() (*stats.Config, error) {
	return &stats.Config{}, nil
}

type Config struct {
	Port            uint16                    `json:"port"` // Port of this Point server.
	LogConfig       *LogConfig                `json:"log"`
	RouterConfig    *RouterConfig             `json:"routing"`
	DNSConfig       *DnsConfig                `json:"dns"`
	InboundConfig   *InboundConnectionConfig  `json:"inbound"`
	OutboundConfig  *OutboundConnectionConfig `json:"outbound"`
	InboundDetours  []InboundDetourConfig     `json:"inboundDetour"`
	OutboundDetours []OutboundDetourConfig    `json:"outboundDetour"`
	Transport       *TransportConfig          `json:"transport"`
	Policy          *PolicyConfig             `json:"policy"`
	Api             *ApiConfig                `json:"api"`
	Stats           *StatsConfig              `json:"stats"`
}

// Build implements Buildable.
func (c *Config) Build() (*core.Config, error) {
	config := &core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&dispatcher.Config{}),
			serial.ToTypedMessage(&proxyman.InboundConfig{}),
			serial.ToTypedMessage(&proxyman.OutboundConfig{}),
		},
	}

	if c.Api != nil {
		apiConf, err := c.Api.Build()
		if err != nil {
			return nil, err
		}
		config.App = append(config.App, serial.ToTypedMessage(apiConf))
	}

	if c.Stats != nil {
		statsConf, err := c.Stats.Build()
		if err != nil {
			return nil, err
		}
		config.App = append(config.App, serial.ToTypedMessage(statsConf))
	}

	if c.LogConfig != nil {
		config.App = append(config.App, serial.ToTypedMessage(c.LogConfig.Build()))
	} else {
		config.App = append(config.App, serial.ToTypedMessage(DefaultLogConfig()))
	}

	if c.Transport != nil {
		ts, err := c.Transport.Build()
		if err != nil {
			return nil, err
		}
		config.Transport = ts
	}

	if c.RouterConfig != nil {
		routerConfig, err := c.RouterConfig.Build()
		if err != nil {
			return nil, err
		}
		config.App = append(config.App, serial.ToTypedMessage(routerConfig))
	}

	if c.DNSConfig != nil {
		config.App = append(config.App, serial.ToTypedMessage(c.DNSConfig.Build()))
	}

	if c.Policy != nil {
		pc, err := c.Policy.Build()
		if err != nil {
			return nil, err
		}
		config.App = append(config.App, serial.ToTypedMessage(pc))
	}

	if c.InboundConfig == nil {
		return nil, newError("no inbound config specified")
	}

	if c.InboundConfig.Port == 0 && c.Port > 0 {
		c.InboundConfig.Port = c.Port
	}

	ic, err := c.InboundConfig.Build()
	if err != nil {
		return nil, err
	}
	config.Inbound = append(config.Inbound, ic)

	for _, rawInboundConfig := range c.InboundDetours {
		ic, err := rawInboundConfig.Build()
		if err != nil {
			return nil, err
		}
		config.Inbound = append(config.Inbound, ic)
	}

	if c.OutboundConfig == nil {
		return nil, newError("no outbound config specified")
	}
	oc, err := c.OutboundConfig.Build()
	if err != nil {
		return nil, err
	}
	config.Outbound = append(config.Outbound, oc)

	for _, rawOutboundConfig := range c.OutboundDetours {
		oc, err := rawOutboundConfig.Build()
		if err != nil {
			return nil, err
		}
		config.Outbound = append(config.Outbound, oc)
	}

	return config, nil
}
