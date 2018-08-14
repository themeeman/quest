package commands

import (
	"github.com/bwmarrin/discordgo"
	commands "../../discordcommands"
	"regexp"
	"strings"
	"strconv"
	"math"
	"reflect"
	"database/sql"
	"fmt"
	"sort"
	"bytes"
	"../structures"
	"github.com/fatih/structs"
)

func getOptions(bot *Bot) map[string]*structures.Option {
	g := new(structures.Guild)
	t := reflect.TypeOf(*g)
	options := make(map[string]*structures.Option)
	for _, s := range structs.Names(*g) {
		field, _ := t.FieldByName(s)
		t, ok := field.Tag.Lookup("type")
		if ok {
			name, _ := field.Tag.Lookup("db")
			options[s] = &structures.Option{
				Name: name,
				Type: t,
			}
		}
	}
	return options
}

func (bot *Bot) Set(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string) error {
	options := getOptions(bot)
	fmt.Println(options)
	if args["Option"] == "" {
		names := make([]string, len(options))
		var i int
		for name := range options {
			names[i] = name
			i++
		}
		sort.Strings(names)
		guild := bot.Guilds.Get(commands.MustGetGuildID(session, message))
		var buf bytes.Buffer
		for _, name := range names {
			current := repr(reflect.Indirect(reflect.ValueOf(guild).Elem()).FieldByName(name).Interface(), options[name].Type)
			buf.WriteString(fmt.Sprintf("**%s** - %s\n", name, current))
		}
		_, err := session.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed{
			Description: buf.String(),
		})
		if err != nil {
			fmt.Println(err)
		}
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

func repr(val interface{}, T string) string {
	switch v := val.(type) {
	case sql.NullString:
		if v.Valid {
			switch T {
			case "ChannelMention":
				return "<#" + v.String + ">"
			case "UserMention":
				return "<@" + v.String + ">"
			case "RoleMention":
				return "<@&" + v.String + ">"
			}
		} else {
			return "None"
		}
	case structures.Modlog:
		if v.Valid {
			return "<#" + v.ChannelID + ">"
		} else {
			return "None"
		}
	default:
		return fmt.Sprint(val)
	}
	return fmt.Sprint(val)
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
	case "ChannelMention":
		if value == "none" {
			a = structures.Modlog{}
		} else if len(value) > 18 {
			a = structures.Modlog{
				ChannelID: value[2:20],
				Valid:  true,
				Log: make(chan structures.Case),
			}
		} else {
			a = structures.Modlog{
				ChannelID: value,
				Valid:  true,
				Log: make(chan structures.Case),
			}
		}
	case "Boolean":
		l := strings.ToLower(value)
		if l == "true" || l == "yes" || l == "y" {
			a = true
		} else {
			a = false
		}
	}
	return a
}
