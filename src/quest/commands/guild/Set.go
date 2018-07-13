package guild

import (
	"github.com/bwmarrin/discordgo"
	commands "../../../discordcommands"
	"regexp"
	"strings"
	"strconv"
	"math"
	"reflect"
	"database/sql"
	"fmt"
	"sort"
	"bytes"
)

func Set(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string, bot *commands.Bot) commands.BotError {
	options := commands.GetOptions(bot)
	fmt.Println(options)
	if args["Option"] == "" {
		names := make([]string, len(options))
		i := 0
		for name := range options {
			names[i] = name
			i++
		}
		sort.Strings(names)
		guild := bot.Guilds.Get(commands.MustGetGuildID(session, message))
		var buf bytes.Buffer
		for _, name := range names {
			current := repr(reflect.Indirect(reflect.ValueOf(guild).Elem()).FieldByName(name).Interface())
			buf.WriteString(fmt.Sprintf("**%s** - %s\n", name, current))
		}
		session.ChannelMessageSend(message.ChannelID, buf.String())
	} else if args["Value"] == "" {
		return commands.UsageError{
			Usage: bot.CommandMap["set"].GetUsage(bot.Prefix, "set"),
		}
	} else {
		keyName := args["Option"]
		option, ok := options[keyName]
		if !ok {
			var found bool
			for k, o := range options {
				if strings.ToLower(args["Option"]) == strings.ToLower(k) {
					option = o
					found = true
					keyName = k
					break
				}
			}
			if !found {
				return commands.OptionError{
					Option: args["Option"],
				}
			}
		}
		pattern, _ := bot.Regex[option.Type]
		value := args["Value"]
		result, _ := regexp.MatchString(pattern, value)
		if !result {
			return commands.ParsingError{
				Value:    value,
				Position: 2,
				Expected: option.Type,
			}
		}
		guild := bot.Guilds.Get(commands.MustGetGuildID(session, message))
		field := reflect.ValueOf(guild).Elem().FieldByName(keyName)
		fieldType := reflect.TypeOf(field.Interface())
		val := reflect.ValueOf(convertType(message, option.Type, value)).Convert(fieldType).Interface()
		field.Set(reflect.ValueOf(val))
		session.MessageReactionAdd(message.ChannelID, message.ID, "☑")
	}
	return nil
}

func repr(val interface{}) string {
	switch val.(type) {
	case sql.NullString:
		if val.(sql.NullString).Valid {
			return val.(sql.NullString).String
		} else {
			return "None"
		}
	default:
		return fmt.Sprint(val)
	}
}

func convertType(message *discordgo.MessageCreate, T string, value string) interface{} {
	var a interface{}
	switch T {
	case "Integer", "SignedInteger":
		a, _ = strconv.Atoi(value)
	case "Decimal", "Float", "BigNumber":
		c := strings.Split(value, "e")
		v, _ := strconv.ParseFloat(c[0], 32)
		if len(c) == 1 {
			a = v
		}
		e, _ := strconv.Atoi(c[1])
		a = v * math.Pow10(e)
		break
	case "UserMention":
		if value == "none" {
			a = sql.NullString{}
		} else if len(message.Mentions) > 0 {
			a = sql.NullString{
				String: message.Mentions[0].ID,
				Valid:  true,
			}
		} else {
			a = sql.NullString{
				String: value,
				Valid:  true,
			}
		}
		break
	case "RoleMention":
		if value == "none" {
			a = sql.NullString{}
		} else if len(message.MentionRoles) > 0 {
			a = sql.NullString{
				String: message.MentionRoles[0],
				Valid:  true,
			}
		} else {
			a = sql.NullString{
				String: value,
				Valid:  true,
			}
		}
		break
	case "ChannelMention":
		if value == "none" {
			a = sql.NullString{}
		} else if len(value) > 18 {
			a = sql.NullString{
				String: value[2:20],
				Valid:  true,
			}
		} else {
			a = sql.NullString{
				String: value,
				Valid:  true,
			}
		}
		break
	case "Boolean":
		l := strings.ToLower(value)
		if l == "true" || l == "yes" || l == "y" {
			a = true
		} else {
			a = false
		}
		break
	}
	return a
}