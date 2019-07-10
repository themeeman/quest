package commands

import (
	"fmt"
	"github.com/pkg/errors"
)

var UserNotFoundError = errors.New(`Could not find user in the server
Maybe they left because its so trash?`)

type UnknownCommandError struct {
	Command string
}

func (e UnknownCommandError) Error() string {
	return fmt.Sprintf(`Unknown command: %s
Use the help command for a list of commands`, e.Command)
}

type MutedError struct {
	Username      string
	Discriminator string
}

func (e MutedError) Error() string {
	return fmt.Sprintf(`Cannot mute user %s#%s
They are already muted!`, e.Username, e.Discriminator)
}

type UnmutedError struct {
	Username      string
	Discriminator string
}

func (e UnmutedError) Error() string {
	return fmt.Sprintf(`Cannot unmute user %s#%s
They aren't muted!`, e.Username, e.Discriminator)
}

type RoleError struct {
	ID string
}

func (e RoleError) Error() string {
	return fmt.Sprintf("Unable to find role '%s'", e.ID)
}
