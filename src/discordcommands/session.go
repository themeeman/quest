package discordcommands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"reflect"
	"runtime/debug"
	"strings"
	"time"
)

func getCommand(commands CommandMap, handlers HandlerMap, name string) (*Command, Handler, string) {
	name = strings.ToLower(name)
	command, okc := commands[name]
	handler, okh := handlers[name]
	if !okc || !okh {
		for n, cmd := range commands {
			for _, alias := range cmd.Aliases {
				if name == alias {
					return getCommand(commands, handlers, n)
				}
			}
		}
	}
	return command, handler, name
}

func NewSession(bot interface{}, token string) (*discordgo.Session, error) {
	t := reflect.TypeOf(bot)
	v := reflect.ValueOf(bot)
	var prefix *string
	if field, ok := t.Elem().FieldByName("Prefix"); ok && field.Type.Kind() == reflect.String {
		prefix = v.Elem().FieldByName("Prefix").Addr().Interface().(*string)
	} else {
		*prefix = ""
	}
	var commands CommandMap
	if field, ok := t.Elem().FieldByName("CommandMap"); ok && field.Type == reflect.TypeOf(commands) {
		commands = v.Elem().FieldByName("CommandMap").Interface().(CommandMap)
	}
	var handlers = make(HandlerMap)
	fmt.Println(t.NumMethod())
	for i := 0; i < t.NumMethod(); i++ {
		funcValue := v.Method(i)
		funcType := v.Method(i).Type()
		handlerType := reflect.TypeOf(Handler(nil))
		if funcType.ConvertibleTo(handlerType) {
			fmt.Println(strings.ToLower(t.Method(i).Name))
			handlers[strings.ToLower(t.Method(i).Name)] = funcValue.Convert(handlerType).Interface().(Handler)
		}
	}
	var types map[string]string
	if field, ok := t.Elem().FieldByName("Types"); ok && field.Type == reflect.TypeOf(types) {
		types = v.Elem().FieldByName("Types").Interface().(map[string]string)
	}
	var errors chan struct {
		Err error
		*discordgo.MessageCreate
	}
	if field, ok := t.Elem().FieldByName("Errors"); ok && field.Type == reflect.TypeOf(errors) {
		errors = v.Elem().FieldByName("Errors").Interface().(chan struct {
			Err error
			*discordgo.MessageCreate
		})
	}
	var groupNames map[Group]string
	if field, ok := t.Elem().FieldByName("GroupNames"); ok && field.Type == reflect.TypeOf(groupNames) {
		groupNames = v.Elem().FieldByName("GroupNames").Interface().(map[Group]string)
	}
	var userGroup func(session *discordgo.Session, guild *discordgo.Guild, member *discordgo.Member) Group
	if _, ok := t.MethodByName("UserGroup"); ok {
		if f, ok := v.MethodByName("UserGroup").Interface().(func(session *discordgo.Session, guild *discordgo.Guild, member *discordgo.Member) Group); ok {
			userGroup = f
		}
	}
	var execute func(*discordgo.Session, *discordgo.MessageCreate)
	execute = func(session *discordgo.Session, message *discordgo.MessageCreate) {
		defer func() {
			if r := recover(); r != nil {
				log.Println(string(debug.Stack()))
				session.ChannelMessageSend(message.ChannelID, "```"+`An unexpected panic occured in the execution of that command.
`+fmt.Sprint(r)+"\nTry again later, or contact themeeman#8354"+"```")
			}
		}()
		t := time.Now()
		if !HasPrefix(message.Content, *prefix) {
			return
		}
		args := strings.Fields(TrimPrefix(message.Content, *prefix))
		if len(args) == 0 {
			return
		}
		info, cmd, name := getCommand(commands, handlers, args[0])
		if cmd == nil {
			return
		}
		g, _ := session.Guild(MustGetGuildID(session, message))
		if g == nil {
			return
		}
		m, _ := session.GuildMember(g.ID, message.Author.ID)
		if m == nil {
			return
		}
		if userGroup != nil {
			if level := userGroup(session, g, m); level < info.Group {
				if errors != nil {
					errors <- struct {
						Err error
						*discordgo.MessageCreate
					}{
						Err: InsufficientPermissionsError{
							Required: groupNames[info.Group],
							Had:      groupNames[level],
						},
						MessageCreate: message,
					}
				}
				return
			}
		}
		args = args[1:]
		if len(args) == 0 && info != nil && info.ForcedArgs() > 0 {
			if errors != nil {
				errors <- struct {
					Err error
					*discordgo.MessageCreate
				}{
					Err:           ZeroArgumentsError{Command: name},
					MessageCreate: message,
				}
			}
			return
		}
		var newArgs = map[string]string{}
		if types != nil {
			var err error
			newArgs, err = parseArgs(types, info, args)
			if err != nil {
				if errors != nil {
					if e, ok := err.(UsageError); ok {
						e.Usage = info.GetUsage(*prefix, name)
						err = e
					}
					errors <- struct {
						Err error
						*discordgo.MessageCreate
					}{
						Err:           err,
						MessageCreate: message,
					}
				}
				return
			}
		}
		fmt.Println(newArgs)
		err := cmd(session, message, newArgs)
		if err != nil {
			if errors != nil {
				errors <- struct {
					Err error
					*discordgo.MessageCreate
				}{
					Err:           err,
					MessageCreate: message,
				}
			}
			return
		}
		fmt.Println(time.Since(t))
		return
	}
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	dg.AddHandler(execute)
	return dg, nil
}
