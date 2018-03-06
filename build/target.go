package build

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"v2ray.com/ext/gpg"
	"v2ray.com/ext/sysio"
	"v2ray.com/ext/zip"
)

type Output struct {
	Generated string
}

type Target interface {
	BuildTo(directory string) (*Output, error)
}

type ResourceTarget struct {
	Source           Path
	Target           string
	OS               OS
	FixLineSeparator bool
}

func (t *ResourceTarget) BuildTo(dir string) (*Output, error) {
	content, err := sysio.ReadFile(t.Source.Get())
	if err != nil {
		return nil, err
	}
	if t.FixLineSeparator {
		content = bytes.Replace(content, []byte{'\r', '\n'}, []byte{'\n'}, -1)
		content = bytes.Replace(content, []byte{'\r'}, []byte{'\n'}, -1)
		switch t.OS {
		case Windows:
			content = bytes.Replace(content, []byte{'\n'}, []byte{'\r', '\n'}, -1)
		case MacOS:
			content = bytes.Replace(content, []byte{'\n'}, []byte{'\r'}, -1)
		}
	}

	target := filepath.Join(dir, t.Target)
	if err := os.MkdirAll(filepath.Dir(target), os.ModePerm); err != nil {
		return nil, err
	}
	if err := ioutil.WriteFile(target, content, 0777); err != nil {
		return nil, err
	}
	return &Output{Generated: target}, nil
}

type ZipTarget struct {
	Options []zip.Option
	Target  string
	Source  Path
}

func (zo *ZipTarget) BuildTo(dir string) (*Output, error) {
	target := filepath.Join(dir, zo.Target)
	target = os.ExpandEnv(target)
	src := zo.Source.Get()
	if err := zip.CompressFolder(src, target, zo.Options...); err != nil {
		return nil, err
	}

	return &Output{Generated: target}, nil
}

type SignTarget struct {
	Passphrase string
	Source     Path
}

func (so *SignTarget) BuildTo(dir string) (*Output, error) {
	src := so.Source.Get()
	// TODO: Modify gpg package to accept target file parameter.
	if err := gpg.SignFile(src, so.Passphrase); err != nil {
		return nil, err
	}
	return &Output{Generated: src + ".sig"}, nil
}

type CachedTarget struct {
	Target
	*Output
}

func (t *CachedTarget) BuildTo(dir string) (*Output, error) {
	if t.Output == nil {
		output, err := t.Target.BuildTo(dir)
		if err != nil {
			return nil, err
		}
		t.Output = output
	}
	return t.Output, nil
}
