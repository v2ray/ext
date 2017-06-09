package build

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func BuildV2RayCore(target string, goOS GoOS, goArch GoArch, disableConsoleForWindows bool) error {
	ldFlags := []string{"-s"}
	if disableConsoleForWindows {
		ldFlags = append(ldFlags, "-H windowsgui")
	}
	version := os.Getenv("TRAVIS_TAG")
	if len(version) > 0 {
		year, month, day := time.Now().UTC().Date()
		today := fmt.Sprintf("%04d%02d%02d", year, int(month), day)
		ldFlags = append(ldFlags, " -X v2ray.com/core.version="+version, " -X v2ray.com/core.build="+today)

		bUser := os.Getenv("V_USER")
		if len(bUser) > 0 {
			ldFlags = append(ldFlags, " -X v2ray.com/ext/tools/conf.bUser="+bUser)
		}
	}
	return GoBuild("v2ray.com/core/main", target, goOS, goArch, strings.Join(ldFlags, " "), "json")
}
