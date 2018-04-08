package command

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg command -path Ext,Tools,Conf,Command

import (
	"os"

	"github.com/gogo/protobuf/proto"
	"v2ray.com/core/common"
	"v2ray.com/ext/tools/conf/serial"
	"v2ray.com/ext/tools/control"
)

type ConfigCommand struct{}

func (c *ConfigCommand) Name() string {
	return "config"
}

func (c *ConfigCommand) Description() control.Description {
	return control.Description{
		Short: "Convert config among different formats.",
		Usage: []string{
			"v2ctl config",
		},
	}
}

func (c *ConfigCommand) Execute(args []string) error {
	pbConfig, err := serial.LoadJSONConfig(os.Stdin)
	if err != nil {
		os.Stderr.WriteString("failed to parse json config: " + err.Error())
		os.Exit(-1)
	}

	bytesConfig, err := proto.Marshal(pbConfig)
	if err != nil {
		os.Stderr.WriteString("failed to marshal proto config: " + err.Error())
		os.Exit(-1)
	}

	if _, err := os.Stdout.Write(bytesConfig); err != nil {
		os.Stderr.WriteString("failed to write proto config: " + err.Error())
		os.Exit(-1)
	}
	return nil
}

func init() {
	common.Must(control.RegisterCommand(&ConfigCommand{}))
}
