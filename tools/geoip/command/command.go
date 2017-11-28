package command

import (
	"fmt"

	"v2ray.com/core/common"
	"v2ray.com/ext/tools/control"
)

func init() {
	const name = "geoip"
	common.Must(control.RegisterCommand(name, "Control over GeoIP data.", func(args []string) {
		if len(args) == 0 {
			fmt.Println()
		}
	}))
}
