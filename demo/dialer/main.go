package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"v2ray.com/core"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/proxyman"
	_ "v2ray.com/core/app/proxyman/outbound"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/socks"
	_ "v2ray.com/core/transport/internet/tcp"
)

func main() {
	config := &core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&dispatcher.Config{}),
			serial.ToTypedMessage(&proxyman.OutboundConfig{}),
		},
		Outbound: []*core.OutboundHandlerConfig{{
			ProxySettings: serial.ToTypedMessage(&socks.ClientConfig{
				Server: []*protocol.ServerEndpoint{{
					Address: net.NewIPOrDomain(net.ParseAddress("162.243.108.129")),
					Port:    1080,
				}},
			})},
		},
	}

	v, err := core.New(config)
	if err != nil {
		fmt.Println("Failed to create V: ", err.Error())
		os.Exit(-1)
	}

	conn, err := core.Dial(context.Background(), v, net.TCPDestination(net.ParseAddress("www.v2ray.com"), net.Port(80)))
	if err != nil {
		fmt.Println("Failed to dial connection: ", err.Error())
	}

	_, err = conn.Write([]byte(`GET / HTTP/1.1
Host: www.v2ray.com
User-Agent: Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US; rv:1.9.1.5) Gecko/20091102 Firefox/3.5.5 (.NET CLR 3.5.30729)
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8
Accept-Language: en-us,en;q=0.5
Accept-Encoding: gzip,deflate
Accept-Charset: ISO-8859-1,utf-8;q=0.7,*;q=0.7

`))
	if err != nil {
		fmt.Println("Failed to write request: ", err.Error())
	}

	_, err = io.Copy(os.Stdout, conn)
	if err != nil {
		fmt.Println("Failed to read response: ", err.Error())
	}
}
