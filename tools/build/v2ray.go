package build

import (
	"fmt"
	"os"
	"time"
)

var releases = map[string][]*GoBuildTarget{
	genReleaseId(Windows, Amd64):  append(genRegularTarget(Windows, Amd64), getWindowsExtra(Amd64)...),
	genReleaseId(Windows, X86):    append(genRegularTarget(Windows, X86), getWindowsExtra(X86)...),
	genReleaseId(MacOS, Amd64):    append(genRegularTarget(MacOS, Amd64)),
	genReleaseId(Linux, Amd64):    append(genRegularTarget(Linux, Amd64)),
	genReleaseId(Linux, X86):      append(genRegularTarget(Linux, X86)),
	genReleaseId(Linux, Arm64):    append(genRegularTarget(Linux, Arm64)),
	genReleaseId(Linux, Arm):      append(genRegularTarget(Linux, Arm), getArmExtra()...),
	genReleaseId(Linux, Mips64):   append(genRegularTarget(Linux, Mips64)),
	genReleaseId(Linux, Mips64LE): append(genRegularTarget(Linux, Mips64LE)),
	genReleaseId(Linux, Mips):     append(genRegularTarget(Linux, Mips), getMipsExtra(Mips)...),
	genReleaseId(Linux, MipsLE):   append(genRegularTarget(Linux, MipsLE), getMipsExtra(MipsLE)...),
	genReleaseId(Linux, S390X):    append(genRegularTarget(Linux, S390X)),
	genReleaseId(OpenBSD, Amd64):  append(genRegularTarget(OpenBSD, Amd64)),
	genReleaseId(OpenBSD, X86):    append(genRegularTarget(OpenBSD, X86)),
	genReleaseId(FreeBSD, Amd64):  append(genRegularTarget(FreeBSD, Amd64)),
	genReleaseId(FreeBSD, X86):    append(genRegularTarget(FreeBSD, X86)),
}

func genReleaseId(goOS GoOS, goArch GoArch) string {
	return string(goOS) + "-" + string(goArch)
}

const stdSource = "v2ray.com/core/main"
const stdTarget = "v2ray"
const stdControlSource = "v2ray.com/ext/tools/control/main"
const stdControlTarget = "v2ctl"

func targetWithSuffix(goOS GoOS, target string) string {
	if goOS == Windows {
		return target + ".exe"
	}
	return target
}

func genRegularTarget(goOS GoOS, goArch GoArch) []*GoBuildTarget {
	return []*GoBuildTarget{
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

func getWindowsExtra(goArch GoArch) []*GoBuildTarget {
	return []*GoBuildTarget{
		{
			Source:  stdSource,
			Target:  "w" + stdTarget + ".exe",
			OS:      Windows,
			Arch:    goArch,
			LdFlags: append(genStdLdFlags(Windows, goArch), "-H windowsgui"),
		},
	}
}

func getArmExtra() []*GoBuildTarget {
	return []*GoBuildTarget{
		{
			Source:  stdSource,
			Target:  stdTarget + "_armv7",
			OS:      Linux,
			Arch:    Arm,
			LdFlags: genStdLdFlags(Linux, Arm),
			ArmOpt:  "7",
		},
	}
}

func getMipsExtra(goArch GoArch) []*GoBuildTarget {
	return []*GoBuildTarget{
		{
			Source:  stdSource,
			Target:  stdTarget + "_softfloat",
			OS:      Linux,
			Arch:    goArch,
			LdFlags: genStdLdFlags(Linux, goArch),
			MipsOpt: "softfloat",
		},
		{
			Source:  stdControlSource,
			Target:  stdControlTarget + "_softfloat",
			OS:      Linux,
			Arch:    goArch,
			LdFlags: genStdLdFlags(Linux, goArch),
			MipsOpt: "softfloat",
		},
	}
}

func genStdLdFlags(goOS GoOS, goArch GoArch) []string {
	ldFlags := []string{"-s -w"}
	version := os.Getenv("TRAVIS_TAG")
	if len(version) > 0 {
		year, month, day := time.Now().UTC().Date()
		today := fmt.Sprintf("%04d%02d%02d", year, int(month), day)
		ldFlags = append(ldFlags, " -X v2ray.com/core.version="+version, " -X v2ray.com/core.build="+today)
	}
	return ldFlags
}

func GetReleaseTargets(goOS GoOS, goArch GoArch) []*GoBuildTarget {
	return releases[genReleaseId(goOS, goArch)]
}
