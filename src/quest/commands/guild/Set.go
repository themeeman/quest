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
			current := repr(reflect.Indirect(reflect.ValueOf(*guild)).FieldByName(name).Interface())
			buf.WriteString(fmt.Sprintf("**%s** - %s\n", name, current))
		}
		session.ChannelMessageSend(message.ChannelID, buf.String())
	} else if args["Value"] == "" {
		return commands.InsufficentArgumentsError{
			Minimum:  2,
			Received: 1,
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
		val := convertType(message, option.Type, value)
		guild := bot.Guilds.Get(commands.MustGetGuildID(session, message))
		field := reflect.ValueOf(guild).Elem().FieldByName(keyName)
		if n, ok := val.(int); ok {
			val = convertInt(n, field)
		} else if n, ok := val.(float64); ok {
			val = convertFloat(n, field)
		} else if n, ok := val.(complex128); ok {
			val = convertComplex(n, field)
		}
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


func convertInt(val int, field reflect.Value) interface{} {
	var a interface{}
	switch field.Interface().(type) {
	case uint8: a = uint8(val)
	case uint16: a = uint16(val)
	case uint32: a = uint32(val)
	case uint64: a = uint64(val)
	case int8: a = int8(val)
	case int16: a = int16(val)
	case int32: a = int32(val)
	case int64: a = int64(val)
	case int: a = val
	case uint: a = uint(val)
	}
	return a
}

func convertFloat(val float64, field reflect.Value) interface{} {
	var a interface{}
	switch field.Interface().(type) {
	case float32: a = float32(val)
	case float64: a = float64(val)
	}
	return a
}

func convertComplex(val complex128, field reflect.Value) interface{} {
	var a interface{}
	switch field.Interface().(type) {
	case complex64: a = complex64(val)
	case complex128: a = complex128(val)
	}
	return a
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