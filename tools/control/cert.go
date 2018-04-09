package control

import (
	"encoding/json"
	"flag"
	"os"
	"strings"

	"v2ray.com/core/common"
	"v2ray.com/core/common/protocol/tls/cert"
)

type stringList []string

func (l *stringList) String() string {
	return "String list"
}

func (l *stringList) Set(v string) error {
	if len(v) == 0 {
		return newError("empty value")
	}
	*l = append(*l, v)
	return nil
}

type jsonCert struct {
	Certificate []string `json:"certificate"`
	Key         []string `json:"key"`
}

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

func (c *CertificateCommand) printJson(certificate *cert.Certificate) {
	jCert := &jsonCert{
		Certificate: strings.Split(strings.TrimSpace(string(certificate.Certificate)), "\n"),
		Key:         strings.Split(strings.TrimSpace(string(certificate.PrivateKey)), "\n"),
	}
	content, err := json.MarshalIndent(jCert, "", "  ")
	common.Must(err)
	os.Stdout.Write(content)
}

func (c *CertificateCommand) Execute(args []string) error {
	fs := flag.NewFlagSet(c.Name(), flag.ContinueOnError)

	var domainNames stringList
	fs.Var(&domainNames, "domain", "Domain name for the certificate")

	isCA := fs.Bool("ca", false, "Whether this certificate is a CA")
	output := fs.String("output", "json", "Output format")

	var opts []cert.Option
	if *isCA {
		opts = append(opts, cert.Authority(*isCA))
	}

	cert, err := cert.Generate(nil, opts...)
	if err != nil {
		return newError("failed to generate TLS certificate")
	}

	switch strings.ToLower(*output) {
	case "json":
		c.printJson(cert)
	case "file":
	default:
		return newError("unknown output format: ", *output)
	}

	return nil
}

func init() {
	common.Must(RegisterCommand(&CertificateCommand{}))
}
