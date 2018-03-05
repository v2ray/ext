package build

import (
	"fmt"
	"os"
	"time"

	"v2ray.com/ext/build"
)

type releaseID struct {
	OS   build.OS
	Arch build.Arch
}

var releases = map[releaseID][]*build.GoTarget{
	genReleaseID(build.Windows, build.Amd64):  append(genRegularTarget(build.Windows, build.Amd64), getWindowsExtra(build.Amd64)...),
	genReleaseID(build.Windows, build.X86):    append(genRegularTarget(build.Windows, build.X86), getWindowsExtra(build.X86)...),
	genReleaseID(build.MacOS, build.Amd64):    append(genRegularTarget(build.MacOS, build.Amd64)),
	genReleaseID(build.Linux, build.Amd64):    append(genRegularTarget(build.Linux, build.Amd64)),
	genReleaseID(build.Linux, build.X86):      append(genRegularTarget(build.Linux, build.X86)),
	genReleaseID(build.Linux, build.Arm64):    append(genRegularTarget(build.Linux, build.Arm64)),
	genReleaseID(build.Linux, build.Arm):      append(genRegularTarget(build.Linux, build.Arm), getArmExtra()...),
	genReleaseID(build.Linux, build.Mips64):   append(genRegularTarget(build.Linux, build.Mips64)),
	genReleaseID(build.Linux, build.Mips64LE): append(genRegularTarget(build.Linux, build.Mips64LE)),
	genReleaseID(build.Linux, build.Mips):     append(genRegularTarget(build.Linux, build.Mips), getMipsExtra(build.Mips)...),
	genReleaseID(build.Linux, build.MipsLE):   append(genRegularTarget(build.Linux, build.MipsLE), getMipsExtra(build.MipsLE)...),
	genReleaseID(build.Linux, build.S390X):    append(genRegularTarget(build.Linux, build.S390X)),
	genReleaseID(build.OpenBSD, build.Amd64):  append(genRegularTarget(build.OpenBSD, build.Amd64)),
	genReleaseID(build.OpenBSD, build.X86):    append(genRegularTarget(build.OpenBSD, build.X86)),
	genReleaseID(build.FreeBSD, build.Amd64):  append(genRegularTarget(build.FreeBSD, build.Amd64)),
	genReleaseID(build.FreeBSD, build.X86):    append(genRegularTarget(build.FreeBSD, build.X86)),
}

func genReleaseID(goOS build.OS, goArch build.Arch) releaseID {
	return releaseID{OS: goOS, Arch: goArch}
}

const stdSource = "v2ray.com/core/main"
const stdTarget = "v2ray"
const stdControlSource = "v2ray.com/ext/tools/control/main"
const stdControlTarget = "v2ctl"

func targetWithSuffix(goOS build.OS, target string) string {
	if goOS == build.Windows {
		return target + ".exe"
	}
	return target
}

func genRegularTarget(goOS build.OS, goArch build.Arch) []*build.GoTarget {
	return []*build.GoTarget{
		{
			Source:  stdSource,
			Target:  targetWithSuffix(goOS, stdTarget),
			OS:      goOS,
			Arch:    goArch,
			LdFlags: genStdLdFlags(goOS, goArch),
		},
		{
			Source:  stdControlSource,
			Target:  targetWithSuffix(goOS, stdControlTarget),
			OS:      goOS,
			Arch:    goArch,
			LdFlags: []string{"-s", "-w"},
		},
	}
}

func getWindowsExtra(goArch build.Arch) []*build.GoTarget {
	return []*build.GoTarget{
		{
			Source:  stdSource,
			Target:  "w" + stdTarget + ".exe",
			OS:      build.Windows,
			Arch:    goArch,
			LdFlags: append(genStdLdFlags(build.Windows, goArch), "-H windowsgui"),
		},
	}
}

func getArmExtra() []*build.GoTarget {
	return []*build.GoTarget{
		{
			Source:  stdSource,
			Target:  stdTarget + "_armv7",
			OS:      build.Linux,
			Arch:    build.Arm,
			LdFlags: genStdLdFlags(build.Linux, build.Arm),
			ArmOpt:  "7",
		},
	}
}

func getMipsExtra(goArch build.Arch) []*build.GoTarget {
	return []*build.GoTarget{
		{
			Source:  stdSource,
			Target:  stdTarget + "_softfloat",
			OS:      build.Linux,
			Arch:    goArch,
			LdFlags: genStdLdFlags(build.Linux, goArch),
			MipsOpt: "softfloat",
		},
		{
			Source:  stdControlSource,
			Target:  stdControlTarget + "_softfloat",
			OS:      build.Linux,
			Arch:    goArch,
			LdFlags: genStdLdFlags(build.Linux, goArch),
			MipsOpt: "softfloat",
		},
	}
}

func genStdLdFlags(goOS build.OS, goArch build.Arch) []string {
	ldFlags := []string{"-s -w"}
	version := os.Getenv("TRAVIS_TAG")
	if len(version) > 0 {
		year, month, day := time.Now().UTC().Date()
		today := fmt.Sprintf("%04d%02d%02d", year, int(month), day)
		ldFlags = append(ldFlags, " -X v2ray.com/core.version="+version, " -X v2ray.com/core.build="+today)
	}
	return ldFlags
}

func GetReleaseTargets(goOS build.OS, goArch build.Arch) []*build.GoTarget {
	return releases[genReleaseID(goOS, goArch)]
}
