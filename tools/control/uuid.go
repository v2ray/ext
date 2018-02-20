package control

import (
	"fmt"

	"v2ray.com/core/common"
	"v2ray.com/core/common/uuid"
)

func init() {
	const name = "uuid"
	common.Must(RegisterCommand(name, "Generate new UUIDs", func(arg []string) {
		u := uuid.New()
		fmt.Println(u.String())
	}))
}
