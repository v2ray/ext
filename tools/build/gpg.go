package build

import (
	"os"
	"os/exec"
)

func GPGSignFile(file string, passphrase string) error {
	cmd := exec.Command("gpg", "--digest-algo", "SHA512", "--no-tty", "--batch", "--passphrase", passphrase, "--output", file+".sig", "--detach-sig", file)
	cmd.Env = append(cmd.Env, os.Environ()...)
	output, err := cmd.CombinedOutput()
	if len(output) > 0 {
		os.Stdout.Write(output)
	}
	return err
}

func GPGImport(file string) error {
	cmd := exec.Command("gpg", "--import", file)
	cmd.Env = append(cmd.Env, os.Environ()...)
	output, err := cmd.CombinedOutput()
	if len(output) > 0 {
		os.Stdout.Write(output)
	}
	return err
}
