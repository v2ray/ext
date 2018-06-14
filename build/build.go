package build

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// OS is a GoOS value for target operating system.
type OS string

const (
	Windows      = OS("windows")
	MacOS        = OS("darwin")
	Linux        = OS("linux")
	FreeBSD      = OS("freebsd")
	OpenBSD      = OS("openbsd")
	DragonflyBSD = OS("dragonfly")
	UnknownOS    = OS("unknown")
)

// Arch is a GoArch value for CPU architecture.
type Arch string

const (
	X86         = Arch("386")
	Amd64       = Arch("amd64")
	Arm         = Arch("arm")
	Arm64       = Arch("arm64")
	Mips64      = Arch("mips64")
	Mips64LE    = Arch("mips64le")
	Mips        = Arch("mips")
	MipsLE      = Arch("mipsle")
	S390X       = Arch("s390x")
	UnknownArch = Arch("unknown")
)

func ParseOS(rawOS string) OS {
	osStr := strings.ToLower(rawOS)
	switch osStr {
	case "windows", "win":
		return Windows
	case "darwin", "mac", "macos", "osx":
		return MacOS
	case "linux", "debian", "ubuntu", "redhat", "centos":
		return Linux
	case "freebsd":
		return FreeBSD
	case "openbsd":
		return OpenBSD
	case "dragonfly", "dragonflybsd":
		return DragonflyBSD
	default:
		return UnknownOS
	}
}

func ParseArch(rawArch string) Arch {
	archStr := strings.ToLower(rawArch)
	switch archStr {
	case "x86", "386", "i386":
		return X86
	case "amd64", "x86-64", "x64":
		return Amd64
	case "arm":
		return Arm
	case "arm64":
		return Arm64
	case "mips":
		return Mips
	case "mipsle":
		return MipsLE
	case "mips64":
		return Mips64
	case "mips64le":
		return Mips64LE
	case "s390x":
		return S390X
	default:
		return UnknownArch
	}
}

func GetSuffix(os OS, arch Arch) string {
	suffix := "-custom"
	switch os {
	case Windows:
		switch arch {
		case X86:
			suffix = "-windows-32"
		case Amd64:
			suffix = "-windows-64"
		}
	case MacOS:
		suffix = "-macos"
	case Linux:
		switch arch {
		case X86:
			suffix = "-linux-32"
		case Amd64:
			suffix = "-linux-64"
		case Arm:
			suffix = "-linux-arm"
		case Arm64:
			suffix = "-linux-arm64"
		case Mips64:
			suffix = "-linux-mips64"
		case Mips64LE:
			suffix = "-linux-mips64le"
		case Mips:
			suffix = "-linux-mips"
		case MipsLE:
			suffix = "-linux-mipsle"
		case S390X:
			suffix = "-linux-s390x"
		}
	case FreeBSD:
		switch arch {
		case X86:
			suffix = "-freebsd-32"
		case Amd64:
			suffix = "-freebsd-64"
		case Arm:
			suffix = "-freebsd-arm"
		}
	case OpenBSD:
		switch arch {
		case X86:
			suffix = "-openbsd-32"
		case Amd64:
			suffix = "-openbsd-64"
		}
	case DragonflyBSD:
		switch arch {
		case X86:
			suffix = "-dragonfly-32"
		case Amd64:
			suffix = "-dragonfly-64"
		}
	}

	return suffix
}

func createDirectoryFor(file string) error {
	return os.MkdirAll(filepath.Dir(file), os.ModePerm)
}

type GoTarget struct {
	Source  string
	Target  string
	OS      OS
	Arch    Arch
	LdFlags []string
	ArmOpt  string
	MipsOpt string
	Tags    []string
}

// Envs returns the environment variables for this build.
func (t *GoTarget) Envs() []string {
	envs := []string{"GOOS=" + string(t.OS), "GOARCH=" + string(t.Arch), "CGO_ENABLED=0"}
	if len(t.ArmOpt) > 0 {
		envs = append(envs, "GOARM="+t.ArmOpt)
	}
	if len(t.MipsOpt) > 0 {
		envs = append(envs, "GOMIPS="+t.MipsOpt)
	}
	return envs
}

func (t *GoTarget) BuildTo(directory string) (*Output, error) {
	goPath := os.Getenv("GOPATH")
	targetFile := filepath.Join(directory, t.Target)
	if err := createDirectoryFor(targetFile); err != nil {
		return nil, err
	}
	args := []string{
		"build",
		"-o", targetFile,
		"-compiler", "gc",
		"-gcflags", "-trimpath=" + goPath,
		"-asmflags", "-trimpath=" + goPath,
	}
	if len(t.LdFlags) > 0 {
		args = append(args, "-ldflags", strings.Join(t.LdFlags, " "))
	}
	if len(t.Tags) > 0 {
		args = append(args, "-tags", strings.Join(t.Tags, ","))
	}
	args = append(args, t.Source)

	cmd := exec.Command("go", args...)
	cmd.Env = append(cmd.Env, t.Envs()...)
	cmd.Env = append(cmd.Env, os.Environ()...)
	output, err := cmd.CombinedOutput()
	if len(output) > 0 {
		os.Stdout.Write(output)
	}

	return &Output{Generated: targetFile}, err
}
