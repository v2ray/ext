package build

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"v2ray.com/ext/sysio"
)

type CopyOption func([]byte) []byte

func FormatLineEnding(goOS GoOS) CopyOption {
	if goOS == Windows {
		return func(content []byte) []byte {
			content = bytes.Replace(content, []byte{'\r', '\n'}, []byte{'\n'}, -1)
			content = bytes.Replace(content, []byte{'\n'}, []byte{'\r', '\n'}, -1)
			return content
		}
	}

	return func(content []byte) []byte {
		return bytes.Replace(content, []byte{'\r', '\n'}, []byte{'\n'}, -1)
	}
}

func CopyFile(src string, dest string, options ...CopyOption) error {
	content, err := sysio.ReadFile(src)
	if err != nil {
		return err
	}
	for _, option := range options {
		content = option(content)
	}
	return ioutil.WriteFile(dest, content, 0777)
}

func CopyAllConfigFiles(destDir string, goOS GoOS) error {
	GOPATH := os.Getenv("GOPATH")
	srcDir := filepath.Join(GOPATH, "src", "v2ray.com", "core", "tools", "release", "config")
	src := filepath.Join(srcDir, "vpoint_socks_vmess.json")
	dest := filepath.Join(destDir, "vpoint_socks_vmess.json")
	if goOS == Windows || goOS == MacOS {
		dest = filepath.Join(destDir, "config.json")
	}
	option := FormatLineEnding(goOS)
	if err := CopyFile(src, dest, option); err != nil {
		return err
	}

	src = filepath.Join(srcDir, "geoip.dat")
	dest = filepath.Join(destDir, "geoip.dat")

	if err := CopyFile(src, dest); err != nil {
		return err
	}

	src = filepath.Join(srcDir, "geosite.dat")
	dest = filepath.Join(destDir, "geosite.dat")

	if err := CopyFile(src, dest); err != nil {
		return err
	}

	src = filepath.Join(GOPATH, "src", "v2ray.com", "core", "tools", "release", "doc", "readme.md")
	dest = filepath.Join(destDir, "readme.md")

	if err := CopyFile(src, dest, option); err != nil {
		return err
	}

	if goOS == Windows || goOS == MacOS {
		return nil
	}

	src = filepath.Join(srcDir, "vpoint_vmess_freedom.json")
	dest = filepath.Join(destDir, "vpoint_vmess_freedom.json")

	if err := CopyFile(src, dest, option); err != nil {
		return err
	}

	if goOS == Linux {
		if err := os.MkdirAll(filepath.Join(destDir, "systemv"), os.ModeDir|0777); err != nil {
			return err
		}
		src = filepath.Join(srcDir, "systemv", "v2ray")
		dest = filepath.Join(destDir, "systemv", "v2ray")
		if err := CopyFile(src, dest); err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Join(destDir, "systemd"), os.ModeDir|0777); err != nil {
			return err
		}
		src = filepath.Join(srcDir, "systemd", "v2ray.service")
		dest = filepath.Join(destDir, "systemd", "v2ray.service")
		if err := CopyFile(src, dest); err != nil {
			return err
		}
	}

	return nil
}
