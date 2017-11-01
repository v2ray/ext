package control

type Command func(args []string)

var (
	commandRegistry = make(map[string]Command)
)

func RegisterCommand(name string, cmd Command) error {
	commandRegistry[name] = cmd
	return nil
}

func GetCommand(name string) Command {
	return commandRegistry[name]
}
