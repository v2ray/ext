package build

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"v2ray.com/ext/build"
)

var (
	OptionSign bool = false
)

type releaseID struct {
	OS   build.OS
	Arch build.Arch
}

var releases map[releaseID][]build.Target

func getReleases() map[releaseID][]build.Target {
	if releases == nil {
		releases = map[releaseID][]build.Target{
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
	}
	return releases
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

func targetWithSignature(target build.Target) []build.Target {
	if !OptionSign {
		return []build.Target{target}
	}

	ct := &build.CachedTarget{
		Target: target,
	}
	st := &build.SignTarget{
		Passphrase: os.Getenv("GPG_SIGN_PASS"),
		Source: &build.LazyPath{
			CachedTarget: ct,
		},
	}
	return []build.Target{ct, st}
}

func genRegularTarget(goOS build.OS, goArch build.Arch) []build.Target {
	releaseDir := filepath.Join("${GOPATH}", "src", "v2ray.com", "core", "release")
	var targets []build.Target

	targets = append(targets, targetWithSignature(&build.GoTarget{
		Source:  stdSource,
		Target:  targetWithSuffix(goOS, stdTarget),
		OS:      goOS,
		Arch:    goArch,
		LdFlags: genStdLdFlags(goOS, goArch),
	})...)

	targets = append(targets, targetWithSignature(&build.GoTarget{
		Source:  stdControlSource,
		Target:  targetWithSuffix(goOS, stdControlTarget),
		OS:      goOS,
		Arch:    goArch,
		LdFlags: []string{"-s", "-w"},
	})...)

	targets = append(targets, &build.ResourceTarget{
		Source: build.EnvPath(filepath.Join(releaseDir, "config", "geoip.dat")),
		Target: "geoip.dat",
	}, &build.ResourceTarget{
		Source: build.EnvPath(filepath.Join(releaseDir, "config", "geosite.dat")),
		Target: "geosite.dat",
	}, &build.ResourceTarget{
		Source:           build.EnvPath(filepath.Join(releaseDir, "doc", "readme.md")),
		Target:           "readme.md",
		FixLineSeparator: true,
		OS:               goOS,
	})

	if goOS == build.Linux {
		targets = append(targets, &build.ResourceTarget{
			Source:           build.EnvPath(filepath.Join(releaseDir, "config", "systemv", "v2ray")),
			Target:           filepath.Join("systemv", "v2ray"),
			FixLineSeparator: true,
			OS:               goOS,
		}, &build.ResourceTarget{
			Source:           build.EnvPath(filepath.Join(releaseDir, "config", "systemd", "v2ray.service")),
			Target:           filepath.Join("systemd", "v2ray.service"),
			FixLineSeparator: true,
			OS:               goOS,
		})
	}

	if goOS == build.Windows || goOS == build.MacOS {
		targets = append(targets, &build.ResourceTarget{
			Source: build.EnvPath(filepath.Join(releaseDir, "config", "vpoint_socks_vmess.json")),
			Target: "config.json",
		})
	} else {
		targets = append(targets, &build.ResourceTarget{
			Source: build.EnvPath(filepath.Join(releaseDir, "config", "vpoint_socks_vmess.json")),
			Target: "vpoint_socks_vmess.json",
		}, &build.ResourceTarget{
			Source: build.EnvPath(filepath.Join(releaseDir, "config", "vpoint_vmess_freedom.json")),
			Target: "vpoint_vmess_freedom.json",
		})
	}

	return targets
}

func getWindowsExtra(goArch build.Arch) []build.Target {
	return targetWithSignature(&build.GoTarget{
		Source:  stdSource,
		Target:  "w" + stdTarget + ".exe",
		OS:      build.Windows,
		Arch:    goArch,
		LdFlags: append(genStdLdFlags(build.Windows, goArch), "-H windowsgui"),
	})
}

func getArmExtra() []build.Target {
	return targetWithSignature(&build.GoTarget{
		Source:  stdSource,
		Target:  stdTarget + "_armv7",
		OS:      build.Linux,
		Arch:    build.Arm,
		LdFlags: genStdLdFlags(build.Linux, build.Arm),
		ArmOpt:  "7",
	})
}

func getMipsExtra(goArch build.Arch) []build.Target {
	var targets []build.Target
	targets = append(targets, targetWithSignature(&build.GoTarget{
		Source:  stdSource,
		Target:  stdTarget + "_softfloat",
		OS:      build.Linux,
		Arch:    goArch,
		LdFlags: genStdLdFlags(build.Linux, goArch),
		MipsOpt: "softfloat",
	})...)
	targets = append(targets, targetWithSignature(&build.GoTarget{
		Source:  stdControlSource,
		Target:  stdControlTarget + "_softfloat",
		OS:      build.Linux,
		Arch:    goArch,
		LdFlags: genStdLdFlags(build.Linux, goArch),
		MipsOpt: "softfloat",
	})...)
	return targets
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

func GetReleaseTargets(goOS build.OS, goArch build.Arch) []build.Target {
	return getReleases()[genReleaseID(goOS, goArch)]
}
