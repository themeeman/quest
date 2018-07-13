package discordcommands

import (
	"github.com/bwmarrin/discordgo"
	"time"
	"fmt"
	"regexp"
	"strings"
)

func TimeToTimestamp(t time.Time) string {
	return t.Format("2006-01-02T15:04:05+00:00")
}

func GetCommand(bot *Bot, name string) (*Command, Handler, string) {
	name = strings.ToLower(name)
	command, okc := bot.CommandMap[name]
	handler, okh := bot.HandlerMap[name]
	if !okc || !okh {
		for n, cmd := range bot.CommandMap {
			for _, alias := range cmd.Aliases {
				if name == alias {
					return GetCommand(bot, n)
				}
			}
		}
	}
	return command, handler, name
}

func ExecuteCommand(session *discordgo.Session, message *discordgo.MessageCreate, bot *Bot) BotError {
	t := time.Now()
	ss := TrimPrefix(message.Content, bot.Prefix)
	args := strings.Fields(ss)
	if len(args) == 0 {
		return nil
	}
	info, cmd, cmdName := GetCommand(bot, args[0])
	args = args[1:]
	if cmd == nil || info == nil {
		return nil
	}
	g, _ := session.Guild(MustGetGuildID(session, message))
	if g == nil {
		return nil
	}
	author, err := session.GuildMember(g.ID, message.Author.ID)
	if err != nil {
		return nil
	}
	sufficent, had, required := SufficentPermissions(session, g, author, bot, info)
	if !sufficent {
		s := []string{"member", "moderator", "admin", "owner"}
		return InsufficentPermissionsError{
			Required: s[required],
			Had:      s[had],
		}
	}
	if len(args) == 0 && info.ForcedArgs() > 0 {
		message.Content = fmt.Sprintf("q:help %s", cmdName)
		return ExecuteCommand(session, message, bot)
	}
	newArgs, err := parseArgs(bot, info, args)
	if err != nil {
		return err
	}
	fmt.Println(newArgs)
	err = cmd(session, message, newArgs, bot)
	if err != nil {
		return err
	}
	guild := bot.Guilds.Get(MustGetGuildID(session, message))
	guild.Members.Get(message.Author.ID)
	fmt.Println(time.Since(t))
	return nil
}

func parseArgs(bot *Bot, command *Command, args []string) (newArgs map[string]string, err BotError) {
	newArgs = make(map[string]string)
	for index, argument := range command.Arguments {
		fmt.Println(index, argument)
		value, err := newArgValue(command, argument, args, index, command.ForcedArgs())
		if err != nil {
			return nil, err
		}
		newArgs[argument.Name] = value
		pattern, ok := bot.Regex[argument.Type]
		if value != "" && ok {
			match, _ := regexp.MatchString(pattern, value)
			if !match {
				return nil, ParsingError{
					Value:    value,
					Position: index,
					Expected: argument.Type,
				}
			}
		}
	}
	return
}

func newArgValue(command *Command, argument *Argument, args []string, index int, forcedArgs int) (string, BotError) {
	if index >= len(args) && !argument.Optional {
		return "", UsageError{
			Usage: command.GetUsage("q:", args[0]),
		}
	} else if index >= len(args) && argument.Optional {
		return "", nil
	} else if argument.Infinite {
		return strings.Join(args[index:], " "), nil
	} else {
		return args[index], nil
	}
}