package discordcommands

import "fmt"

type ParsingError struct {
	Value    string
	Position int
	Expected string
}

func (e ParsingError) Error() string {
	return fmt.Sprintf(`Invalid command arguments
Argument position: %v
Expected Type: %s
Received: %s`, e.Position, e.Expected, e.Value)
}

type UsageError struct {
	Usage string
}

func (e UsageError) Error() string {
	return fmt.Sprintf("The command usage is invalid!\n"+
		"The correct way to use the command is `%s`", e.Usage)
}

type InsufficientPermissionsError struct {
	Required string
	Had      string
}

func (e InsufficientPermissionsError) Error() string {
	return fmt.Sprintf(`You can't execute that command!
You currently are %s, and you need %s`, e.Had, e.Required)
}

type ZeroArgumentsError struct {
	Command string
}

func (e ZeroArgumentsError) Error() string {
	return "No arguments were given"
}