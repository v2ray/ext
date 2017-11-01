package main

import (
	"fmt"
	"os"

	"v2ray.com/ext/tools/control"
)

func main() {
	args := os.Args

	cmd := control.GetCommand(args[1])
	if cmd == nil {
		fmt.Println("Unknown command: ", args[1])
		return
	}

	cmd(args[2:])
}
