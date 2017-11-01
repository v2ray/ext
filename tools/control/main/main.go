package main

import (
	"fmt"
	"os"

	_ "v2ray.com/ext/tools/conf/command"
	"v2ray.com/ext/tools/control"
)

func getCommandName() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}
	return ""
}

func main() {
	name := getCommandName()
	cmd := control.GetCommand(name)
	if cmd == nil {
		fmt.Fprintln(os.Stderr, "Unknown command:", name)
		fmt.Fprintln(os.Stderr)

		fmt.Println("v2ctl <command>")
		fmt.Println("Available commands:")
		control.PrintUsage()
		return
	}

	cmd(os.Args[2:])
}
