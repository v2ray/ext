package control

import "v2ray.com/core/common"

type CertificateCommand struct{}

func (c *CertificateCommand) Name() string {
	return "cert"
}

func (c *CertificateCommand) Description() Description {
	return Description{
		Short: "Generate TLS certificates.",
		Usage: []string{},
	}
}

func (c *CertificateCommand) Execute(args []string) error {

	return nil
}

func init() {
	common.Must(RegisterCommand(&CertificateCommand{}))
}
