package discordcommands

import (
	"github.com/bwmarrin/discordgo"
	"time"
	"fmt"
)

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

type ZeroArgumentsError struct {
	Command string
}

func (e ZeroArgumentsError) Error() string {
	return "No arguments were given"
}

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

type UserNotFoundError struct{}

func (e UserNotFoundError) Error() string {
	return fmt.Sprintf(`Could not find user in the server
Maybe they left because its so trash?`)
}

type UnknownError struct{}

func (e UnknownError) Error() string {
	return fmt.Sprintf(`Sorry! I don't understand this error!`)
}

type TimeError struct {
	Duration int
}

func (e TimeError) Error() string {
	return fmt.Sprintf(`Your given duration of %d is invalid.
Durations must be positive!`, e.Duration)
}

type BotPermissionsError struct{}

func (e BotPermissionsError) Error() string {
	return fmt.Sprintf(`The bot permissions have not been set up correctly`)
}

type UsageError struct {
	Usage string
}

func (e UsageError) Error() string {
	return fmt.Sprintf("The command usage is invalid!\n" +
		"The correct way to use the command is `%s`", e.Usage)
}

type RoleError struct {
	ID string
}

func (e RoleError) Error() string {
	return fmt.Sprintf("Unable to find role '%s'", e.ID)
}

type TypeError struct {
	Name string
}

func (e TypeError) Error() string {
	return fmt.Sprintf(`The provided argument for the Type was incorrect:
%s is **not** a Type.
Use q:types to view all types.`, e.Name)
}

type AutoRoleError struct{}

func (e AutoRoleError) Error() string {
	return "There is no autorole configured. Use q:setautorole to create one."
}

type InsufficentPermissionsError struct {
	Required string
	Had      string
}

func (e InsufficentPermissionsError) Error() string {
	return fmt.Sprintf(`You can't execute that command!
You currently are %s, and you need %s`, e.Had, e.Required)
}

type MuteRoleError struct{}

func (e MuteRoleError) Error() string {
	return "No mute role has been configured for the server! Use q:setmuterole"
}

type OptionError struct {
	Option string
}

func (e OptionError) Error() string {
	return fmt.Sprintf("Option %s does not exist for this guild. Try q:set to list options.", e.Option)
}

type CustomError string

func (e CustomError) Error() string {
	return string(e)
}

func ErrorEmbed(e error) *discordgo.MessageEmbed {
	emb := &discordgo.MessageEmbed{
		Type:        "rich",
		Title:       "An error has occurred",
		Timestamp:   TimeToTimestamp(time.Now().UTC()),
		Color:       0x660000,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Quest Bot",
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Info",
				Value: e.Error(),
			},
		},
	}
	return emb
}
