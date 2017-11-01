package main

import (
	"fmt"
	"os"

	_ "v2ray.com/ext/tools/conf/command"
	"v2ray.com/ext/tools/control"
)

func main() {
	args := os.Args

	cmd := control.GetCommand(args[1])
	if cmd == nil {
		fmt.Fprintln(os.Stderr, "Unknown command:", args[1])
		return
	}

	cmd(args[2:])
}
