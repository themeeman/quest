package discordcommands

import (
	"github.com/bwmarrin/discordgo"
	"time"
	"strings"
	"fmt"
	"reflect"
	"runtime"
	"log"
	"runtime/debug"
)

func getFunctionName(i interface{}) string {
	s := strings.Split(runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name(), ".")
	if len(s) > 0 {
		return s[len(s)-1]
	}
	return ""
}

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
	var prefix string
	if field, ok := t.Elem().FieldByName("Prefix"); ok && field.Type.Kind() == reflect.String {
		prefix = v.Elem().FieldByName("Prefix").String()
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
	var regex map[string]string
	if field, ok := t.Elem().FieldByName("Regex"); ok && field.Type == reflect.TypeOf(regex) {
		regex = v.Elem().FieldByName("Regex").Interface().(map[string]string)
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
	fmt.Println(prefix, commands, handlers, regex, errors, groupNames)
	var execute func(*discordgo.Session, *discordgo.MessageCreate)
	execute = func(session *discordgo.Session, message *discordgo.MessageCreate) {
		defer func() {
			if r := recover(); r != nil {
				log.Println(string(debug.Stack()))
				session.ChannelMessageSend(message.ChannelID, "```"+ `An unexpected panic occured in the execution of that command.
`+ fmt.Sprint(r)+ "\nTry again later, or contact themeeman#8354"+ "```")
			}
		}()
		t := time.Now()
		ss := TrimPrefix(message.Content, prefix)
		args := strings.Fields(ss)
		if len(args) == 0 {
			return
		}
		info, cmd, name := getCommand(commands, handlers, args[0])
		args = args[1:]
		if cmd == nil {
			return
		}
		g, _ := session.Guild(MustGetGuildID(session, message))
		if g == nil {
			return
		}
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
		if regex != nil {
			var err error
			newArgs, err = parseArgs(regex, info, args)
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
		}
		fmt.Println(newArgs)
		err := cmd(session, message, newArgs)
		if err != nil {
			if errors != nil {
				errors <- struct {
					Err       error
					*discordgo.MessageCreate
				}{
					Err:       err,
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