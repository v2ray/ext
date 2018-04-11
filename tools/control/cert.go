package control

import (
	"crypto/x509"
	"encoding/json"
	"flag"
	"os"
	"strings"
	"time"

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
		Usage: []string{
			"v2ctl cert [--ca] [--domain=v2ray.com] [--expire=240h]",
			"Generate new TLS certificate",
			"--ca The new certificate is a CA certificate",
			"--domain Common name for the certificate",
			"--exipre Time until certificate expires. 240h = 10 days.",
		},
	}
}

func (c *CertificateCommand) printJson(certificate *cert.Certificate) {
	certPEM, keyPEM := certificate.ToPEM()
	jCert := &jsonCert{
		Certificate: strings.Split(strings.TrimSpace(string(certPEM)), "\n"),
		Key:         strings.Split(strings.TrimSpace(string(keyPEM)), "\n"),
	}
	content, err := json.MarshalIndent(jCert, "", "  ")
	common.Must(err)
	os.Stdout.Write(content)
	os.Stdout.WriteString("\n")
}

func (c *CertificateCommand) Execute(args []string) error {
	fs := flag.NewFlagSet(c.Name(), flag.ContinueOnError)

	var domainNames stringList
	fs.Var(&domainNames, "domain", "Domain name for the certificate")

	isCA := fs.Bool("ca", false, "Whether this certificate is a CA")
	output := fs.String("output", "json", "Output format")
	expire := fs.Duration("expire", time.Hour*24*90 /* 90 days */, "Time until the certificate expires. Default value 3 months.")

	if err := fs.Parse(args); err != nil {
		return err
	}

	var opts []cert.Option
	if *isCA {
		opts = append(opts, cert.Authority(*isCA))
		opts = append(opts, cert.KeyUsage(x509.KeyUsageCertSign|x509.KeyUsageKeyEncipherment|x509.KeyUsageDigitalSignature))
	}

	opts = append(opts, cert.NotAfter(time.Now().Add(*expire)))

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
