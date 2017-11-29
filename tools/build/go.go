package build

import (
	"os"
	"os/exec"
	"strings"
	"fmt"
)

type GoOS string

const (
	Windows   = GoOS("windows")
	MacOS     = GoOS("darwin")
	Linux     = GoOS("linux")
	FreeBSD   = GoOS("freebsd")
	OpenBSD   = GoOS("openbsd")
	UnknownOS = GoOS("unknown")
)

type GoArch string

const (
	X86         = GoArch("386")
	Amd64       = GoArch("amd64")
	Arm         = GoArch("arm")
	Arm64       = GoArch("arm64")
	Mips64      = GoArch("mips64")
	Mips64LE    = GoArch("mips64le")
	Mips        = GoArch("mips")
	MipsLE      = GoArch("mipsle")
	UnknownArch = GoArch("unknown")
)

func ParseOS(rawOS string) GoOS {
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
	default:
		return UnknownOS
	}
}

func ParseArch(rawArch string) GoArch {
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
	default:
		return UnknownArch
	}
}

func GetSuffix(os GoOS, arch GoArch) string {
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
	}

	return suffix
}

func GoBuild(source string, targetFile string, goOS GoOS, goArch GoArch, ldFlags string, tags ...string) error {
	goPath := os.Getenv("GOPATH")
	args := []string{
		"build",
		"-o", targetFile,
		"-compiler", "gc",
		"-gcflags", "-trimpath=" + goPath,
		"-asmflags", "-trimpath=" + goPath,
	}
	if len(ldFlags) > 0 {
		args = append(args, "-ldflags", ldFlags)
	}
	if len(tags) > 0 {
		args = append(args, "-tags", strings.Join(tags, ","))
	}
	args = append(args, source)

	for index, value := range args {
		fmt.Printf("args[%d]=%d \n", index, value)
	}

	cmd := exec.Command("go", args...)
	cmd.Env = append(cmd.Env, "GOOS="+string(goOS), "GOARCH="+string(goArch), "CGO_ENABLED=0")
	cmd.Env = append(cmd.Env, os.Environ()...)
	output, err := cmd.CombinedOutput()
	if len(output) > 0 {
		os.Stdout.Write(output)
	}

	return err
}
