package conf

import (
	"encoding/json"
	"strings"

	"v2ray.com/core/transport/internet/domainsocket"

	"v2ray.com/core/common/serial"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/http"
	"v2ray.com/core/transport/internet/kcp"
	"v2ray.com/core/transport/internet/tcp"
	"v2ray.com/core/transport/internet/tls"
	"v2ray.com/core/transport/internet/websocket"
	"v2ray.com/ext/sysio"
)

var (
	kcpHeaderLoader = NewJSONConfigLoader(ConfigCreatorCache{
		"none":         func() interface{} { return new(NoOpAuthenticator) },
		"srtp":         func() interface{} { return new(SRTPAuthenticator) },
		"utp":          func() interface{} { return new(UTPAuthenticator) },
		"wechat-video": func() interface{} { return new(WechatVideoAuthenticator) },
	}, "type", "")

	tcpHeaderLoader = NewJSONConfigLoader(ConfigCreatorCache{
		"none": func() interface{} { return new(NoOpConnectionAuthenticator) },
		"http": func() interface{} { return new(HTTPAuthenticator) },
	}, "type", "")
)

type KCPConfig struct {
	Mtu             *uint32         `json:"mtu"`
	Tti             *uint32         `json:"tti"`
	UpCap           *uint32         `json:"uplinkCapacity"`
	DownCap         *uint32         `json:"downlinkCapacity"`
	Congestion      *bool           `json:"congestion"`
	ReadBufferSize  *uint32         `json:"readBufferSize"`
	WriteBufferSize *uint32         `json:"writeBufferSize"`
	HeaderConfig    json.RawMessage `json:"header"`
}

// Build implements Buildable.
func (c *KCPConfig) Build() (*serial.TypedMessage, error) {
	config := new(kcp.Config)

	if c.Mtu != nil {
		mtu := *c.Mtu
		if mtu < 576 || mtu > 1460 {
			return nil, newError("invalid mKCP MTU size: ", mtu).AtError()
		}
		config.Mtu = &kcp.MTU{Value: mtu}
	}
	if c.Tti != nil {
		tti := *c.Tti
		if tti < 10 || tti > 100 {
			return nil, newError("invalid mKCP TTI: ", tti).AtError()
		}
		config.Tti = &kcp.TTI{Value: tti}
	}
	if c.UpCap != nil {
		config.UplinkCapacity = &kcp.UplinkCapacity{Value: *c.UpCap}
	}
	if c.DownCap != nil {
		config.DownlinkCapacity = &kcp.DownlinkCapacity{Value: *c.DownCap}
	}
	if c.Congestion != nil {
		config.Congestion = *c.Congestion
	}
	if c.ReadBufferSize != nil {
		size := *c.ReadBufferSize
		if size > 0 {
			config.ReadBuffer = &kcp.ReadBuffer{Size: size * 1024 * 1024}
		} else {
			config.ReadBuffer = &kcp.ReadBuffer{Size: 512 * 1024}
		}
	}
	if c.WriteBufferSize != nil {
		size := *c.WriteBufferSize
		if size > 0 {
			config.WriteBuffer = &kcp.WriteBuffer{Size: size * 1024 * 1024}
		} else {
			config.WriteBuffer = &kcp.WriteBuffer{Size: 512 * 1024}
		}
	}
	if len(c.HeaderConfig) > 0 {
		headerConfig, _, err := kcpHeaderLoader.Load(c.HeaderConfig)
		if err != nil {
			return nil, newError("invalid mKCP header config.").Base(err).AtError()
		}
		ts, err := headerConfig.(Buildable).Build()
		if err != nil {
			return nil, newError("invalid mKCP header config").Base(err).AtError()
		}
		config.HeaderConfig = ts
	}

	return serial.ToTypedMessage(config), nil
}

type TCPConfig struct {
	HeaderConfig json.RawMessage `json:"header"`
}

// Build implements Buildable.
func (c *TCPConfig) Build() (*serial.TypedMessage, error) {
	config := new(tcp.Config)
	if len(c.HeaderConfig) > 0 {
		headerConfig, _, err := tcpHeaderLoader.Load(c.HeaderConfig)
		if err != nil {
			return nil, newError("invalid TCP header config").Base(err).AtError()
		}
		ts, err := headerConfig.(Buildable).Build()
		if err != nil {
			return nil, newError("invalid TCP header config").Base(err).AtError()
		}
		config.HeaderSettings = ts
	}

	return serial.ToTypedMessage(config), nil
}

type WebSocketConfig struct {
	Path    string            `json:"path"`
	Path2   string            `json:"Path"` // The key was misspelled. For backward compatibility, we have to keep track the old key.
	Headers map[string]string `json:"headers"`
}

// Build implements Buildable.
func (c *WebSocketConfig) Build() (*serial.TypedMessage, error) {
	path := c.Path
	if len(path) == 0 && len(c.Path2) > 0 {
		path = c.Path2
	}
	header := make([]*websocket.Header, 0, 32)
	for key, value := range c.Headers {
		header = append(header, &websocket.Header{
			Key:   key,
			Value: value,
		})
	}

	config := &websocket.Config{
		Path:   path,
		Header: header,
	}
	return serial.ToTypedMessage(config), nil
}

type HTTPConfig struct {
	Host *StringList `json:"host"`
	Path string      `json:"path"`
}

func (c *HTTPConfig) Build() (*serial.TypedMessage, error) {
	config := &http.Config{
		Path: c.Path,
	}
	if c.Host != nil {
		config.Host = []string(*c.Host)
	}
	return serial.ToTypedMessage(config), nil
}

type DomainSocketConfig struct {
	Path     string `json:"path"`
	Abstract bool   `json:"abstract"`
}

func (c *DomainSocketConfig) Build() (*serial.TypedMessage, error) {
	return serial.ToTypedMessage(&domainsocket.Config{
		Path:     c.Path,
		Abstract: c.Abstract,
	}), nil
}

type TLSCertConfig struct {
	CertFile string   `json:"certificateFile"`
	CertStr  []string `json:"certificate"`
	KeyFile  string   `json:"keyFile"`
	KeyStr   []string `json:"key"`
	Usage    string   `json:"usage"`
}

func readFileOrString(f string, s []string) ([]byte, error) {
	if len(f) > 0 {
		return sysio.ReadFile(f)
	}
	if len(s) > 0 {
		return []byte(strings.Join(s, "\n")), nil
	}
	return nil, newError("both file and bytes are empty.")
}

func (c *TLSCertConfig) Build() (*tls.Certificate, error) {
	certificate := new(tls.Certificate)

	cert, err := readFileOrString(c.CertFile, c.CertStr)
	if err != nil {
		return nil, newError("failed to parse certificate").Base(err)
	}
	certificate.Certificate = cert

	if len(c.KeyFile) > 0 || len(c.KeyStr) > 0 {
		key, err := readFileOrString(c.KeyFile, c.KeyStr)
		if err != nil {
			return nil, newError("failed to parse key").Base(err)
		}
		certificate.Key = key
	}

	switch strings.ToLower(c.Usage) {
	case "encipherment":
		certificate.Usage = tls.Certificate_ENCIPHERMENT
	case "verify":
		certificate.Usage = tls.Certificate_AUTHORITY_VERIFY
	case "issue":
		certificate.Usage = tls.Certificate_AUTHORITY_ISSUE
	default:
		certificate.Usage = tls.Certificate_ENCIPHERMENT
	}

	return certificate, nil
}

type TLSConfig struct {
	Insecure   bool             `json:"allowInsecure"`
	Certs      []*TLSCertConfig `json:"certificates"`
	ServerName string           `json:"serverName"`
	ALPN       *StringList      `json:"alpn"`
}

// Build implements Buildable.
func (c *TLSConfig) Build() (*serial.TypedMessage, error) {
	config := new(tls.Config)
	config.Certificate = make([]*tls.Certificate, len(c.Certs))
	for idx, certConf := range c.Certs {
		cert, err := certConf.Build()
		if err != nil {
			return nil, err
		}
		config.Certificate[idx] = cert
	}
	serverName := c.ServerName
	config.AllowInsecure = c.Insecure
	if len(c.ServerName) > 0 {
		config.ServerName = serverName
	}
	if c.ALPN != nil && len(*c.ALPN) > 0 {
		config.NextProtocol = []string(*c.ALPN)
	}
	return serial.ToTypedMessage(config), nil
}

type TransportProtocol string

// Build implements Buildable.
func (p TransportProtocol) Build() (internet.TransportProtocol, error) {
	switch strings.ToLower(string(p)) {
	case "tcp":
		return internet.TransportProtocol_TCP, nil
	case "kcp", "mkcp":
		return internet.TransportProtocol_MKCP, nil
	case "ws", "websocket":
		return internet.TransportProtocol_WebSocket, nil
	case "h2", "http":
		return internet.TransportProtocol_HTTP, nil
	case "ds", "domainsocket":
		return internet.TransportProtocol_DomainSocket, nil
	default:
		return internet.TransportProtocol_TCP, newError("Config: unknown transport protocol: ", p)
	}
}

type StreamConfig struct {
	Network      *TransportProtocol  `json:"network"`
	Security     string              `json:"security"`
	TLSSettings  *TLSConfig          `json:"tlsSettings"`
	TCPSettings  *TCPConfig          `json:"tcpSettings"`
	KCPSettings  *KCPConfig          `json:"kcpSettings"`
	WSSettings   *WebSocketConfig    `json:"wsSettings"`
	HTTPSettings *HTTPConfig         `json:"httpSettings"`
	DSSettings   *DomainSocketConfig `json:"dsSettings"`
}

// Build implements Buildable.
func (c *StreamConfig) Build() (*internet.StreamConfig, error) {
	config := &internet.StreamConfig{
		Protocol: internet.TransportProtocol_TCP,
	}
	if c.Network != nil {
		protocol, err := (*c.Network).Build()
		if err != nil {
			return nil, err
		}
		config.Protocol = protocol
	}
	if strings.ToLower(c.Security) == "tls" {
		tlsSettings := c.TLSSettings
		if tlsSettings == nil {
			tlsSettings = &TLSConfig{}
		}
		ts, err := tlsSettings.Build()
		if err != nil {
			return nil, newError("Failed to build TLS config.").Base(err)
		}
		config.SecuritySettings = append(config.SecuritySettings, ts)
		config.SecurityType = ts.Type
	}
	if c.TCPSettings != nil {
		ts, err := c.TCPSettings.Build()
		if err != nil {
			return nil, newError("Failed to build TCP config.").Base(err)
		}
		config.TransportSettings = append(config.TransportSettings, &internet.TransportConfig{
			Protocol: internet.TransportProtocol_TCP,
			Settings: ts,
		})
	}
	if c.KCPSettings != nil {
		ts, err := c.KCPSettings.Build()
		if err != nil {
			return nil, newError("Failed to build mKCP config.").Base(err)
		}
		config.TransportSettings = append(config.TransportSettings, &internet.TransportConfig{
			Protocol: internet.TransportProtocol_MKCP,
			Settings: ts,
		})
	}
	if c.WSSettings != nil {
		ts, err := c.WSSettings.Build()
		if err != nil {
			return nil, newError("Failed to build WebSocket config.").Base(err)
		}
		config.TransportSettings = append(config.TransportSettings, &internet.TransportConfig{
			Protocol: internet.TransportProtocol_WebSocket,
			Settings: ts,
		})
	}
	if c.HTTPSettings != nil {
		ts, err := c.HTTPSettings.Build()
		if err != nil {
			return nil, newError("Failed to build HTTP config.").Base(err)
		}
		config.TransportSettings = append(config.TransportSettings, &internet.TransportConfig{
			Protocol: internet.TransportProtocol_HTTP,
			Settings: ts,
		})
	}
	if c.DSSettings != nil {
		ds, err := c.DSSettings.Build()
		if err != nil {
			return nil, newError("Failed to build DomainSocket config.").Base(err)
		}
		config.TransportSettings = append(config.TransportSettings, &internet.TransportConfig{
			Protocol: internet.TransportProtocol_DomainSocket,
			Settings: ds,
		})
	}
	return config, nil
}

type ProxyConfig struct {
	Tag string `json:"tag"`
}

// Build implements Buildable.
func (v *ProxyConfig) Build() (*internet.ProxyConfig, error) {
	if len(v.Tag) == 0 {
		return nil, newError("Proxy tag is not set.")
	}
	return &internet.ProxyConfig{
		Tag: v.Tag,
	}, nil
}
