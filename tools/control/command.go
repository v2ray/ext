package control

import "fmt"

type Command func(args []string)

type CommandWithDesc struct {
	desc string
	cmd  Command
}

var (
	commandRegistry = make(map[string]*CommandWithDesc)
)

func RegisterCommand(name string, description string, cmd Command) error {
	commandRegistry[name] = &CommandWithDesc{
		desc: description,
		cmd:  cmd,
	}
	return nil
}

func GetCommand(name string) Command {
	cmd, found := commandRegistry[name]
	if !found {
		return nil
	}
	return cmd.cmd
}

func PrintUsage() {
	for name, desc := range commandRegistry {
		fmt.Println("   ", name, "\t\t\t", desc.desc)
	}
}
