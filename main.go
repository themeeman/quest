package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"syscall"
	"strings"
	"time"
	"strconv"
	"encoding/json"
	commands "discordcommands"
)

var QuestCommands = commands.HandlerMap{
	"help":   Help,
	"mute":   Mute,
	"unmute": Unmute,
	"purge":  Purge,
}
var CommandsData commands.DataMap
var RegexVerifiers = map[string]string{
	"String":         ".",
	"Integer":        "[0-9]",
	"UserMention":    "<@!?[0-9]{18}>",
	"RoleMention":    "<@&[0-9]{18}>",
	"ChannelMention": "<#[0-9]{18}>",
}

var Token string

const (
	prefix      = "q:"
	commandFile = "Commands.json"
)

func ready(s *discordgo.Session, _ *discordgo.Ready) {
	s.UpdateStatus(0, "q:")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println(m.Content)
	if strings.HasPrefix(m.Content, prefix) {
		cmd := strings.TrimPrefix(m.Content, prefix)
		args := strings.Fields(cmd)
		ctx := commands.CommandContext{
			HandlerMap: &QuestCommands,
			DataMap:    &CommandsData,
		}
		err := commands.ExecuteCommand(s, m, args, ctx, RegexVerifiers)
		if err != nil {
			s.ChannelMessageSendEmbed(m.ChannelID, commands.ErrorEmbed(err))
		}
	}
}

func Help(s *discordgo.Session, m *discordgo.MessageCreate, args []string, ctx commands.CommandContext) commands.BotError {
	var data = *ctx.DataMap
	if len(args) == 1 {
		var value string
		for name, v := range data {
			value += fmt.Sprintf("**%v - ** %v\n", name, v.Description)
		}
		fields := []*discordgo.MessageEmbedField{
			{
				Name:  "Commands",
				Value: value,
			},
		}
		s.ChannelMessageSendEmbed(m.ChannelID, QuestEmbed("Help", "", fields))
	} else {
		cmdName := args[1]
		cmdInfo, ok := data[cmdName]
		if !ok {
			return commands.UnknownCommandError{
				Command: cmdName,
			}
		}
		var usage string
		usage = prefix + cmdName + " "
		for _, v := range cmdInfo.Arguments {
			if v.Optional {
				usage += fmt.Sprintf("<%s> ", v.Name)
			} else {
				usage += fmt.Sprintf("[%s] ", v.Name)
			}
		}
		fields := []*discordgo.MessageEmbedField{
			{
				Name:  "Usage",
				Value: "```" + usage + "```",
			},
		}
		s.ChannelMessageSendEmbed(m.ChannelID, QuestEmbed(strings.ToTitle(cmdName), cmdInfo.Description, fields))
	}
	return nil
}

func Mute(s *discordgo.Session, m *discordgo.MessageCreate, args []string, _ commands.CommandContext) commands.BotError {
	if len(args) >= 3 {
		ch, _ := s.State.Channel(m.ChannelID)
		guild, _ := s.Guild(ch.GuildID)
		dur, _ := strconv.Atoi(args[2])
		fmt.Println(m.Mentions)
		var user *discordgo.User
		if len(m.Mentions) == 0 {
			user = new(discordgo.User)
		} else {
			user = m.Mentions[0]
		}
		for _, u := range guild.Members {
			for _, r := range u.Roles {
				if r == "413273250131345409" && u.User.ID == user.ID {
					return commands.MutedError{
						Username:      user.Username,
						Discriminator: user.Discriminator,
					}
				}
			}
		}
		if dur < 0 {
			return commands.TimeError{
				Duration: dur,
			}
		}
		err := s.GuildMemberRoleAdd(ch.GuildID, user.ID, "413273250131345409")
		if err != nil {
			if strings.Contains(err.Error(), "HTTP 403 Forbidden") {
				return commands.PermissionsError{}
			} else {
				return commands.UserNotFoundError{}
			}
		}
		go func() {
			time.Sleep(time.Second * time.Duration(dur))
			s.GuildMemberRoleRemove(ch.GuildID, m.Mentions[0].ID, "413273250131345409")
		}()
		s.ChannelMessageSendEmbed(m.ChannelID, QuestEmbed("Success!", fmt.Sprintf("Successfully muted %s#%s!", user.Username, user.Discriminator), nil))
	}
	return nil
}

func Unmute(s *discordgo.Session, m *discordgo.MessageCreate, args []string, _ commands.CommandContext) commands.BotError {
	ch, err := s.Channel(m.ChannelID)
	if err != nil {
		return nil
	}
	err = s.GuildMemberRoleRemove(ch.GuildID, m.Mentions[0].ID, "413273250131345409")
	if err != nil {
		return nil
	}
	return nil
}

func Purge(s *discordgo.Session, m *discordgo.MessageCreate, args []string, _ commands.CommandContext) commands.BotError {
	i, _ := strconv.Atoi(args[1])
	msgs, err := s.ChannelMessages(m.ChannelID, i, "", "", "")
	if err != nil {
		return nil
	}
	ids := make([]string, i)
	for i, v := range msgs {
		ids[i] = v.ID
		fmt.Println(v.Content)
	}
	//s.ChannelMessagesBulkDelete(m.ChannelID, ids)
	return nil
}

func QuestEmbed(title string, description string, fields []*discordgo.MessageEmbedField) *discordgo.MessageEmbed {
	emb := &discordgo.MessageEmbed{
		Type:      "rich",
		Title:     title,
		Timestamp: ConvertTimeToTimestamp(time.Now()),
		Color:     0x00ffff,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Quest Bot",
		},
		Fields: fields,
	}
	if description != "" {
		emb.Description = description
	}
	return emb
}

func ConvertTimeToTimestamp(t time.Time) string {
	return t.Format("2006-01-02T15:04:05+00:00")
}

func init() {
	func(){
		f, err := os.Open(commandFile)
		if err != nil {
			log.Fatalln("Error opening command file", err)
		}
		stat, err := f.Stat()
		if err != nil {
			log.Fatalln("Error making file stats", err)
		}
		data := make([]byte, stat.Size())
		_, err = f.Read(data)
		if err != nil {
			log.Fatalln("Error reading command file", err)
		}
		err = json.Unmarshal(data, &CommandsData)
		if err != nil {
			log.Fatalln("Error unmarshalling command json", err)
		}
	}()
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalln("Error making discord session", err)
		return
	}

	dg.AddHandler(ready)
	dg.AddHandler(messageCreate)
	err = dg.Open()
	if err != nil {
		log.Fatalln("Error opening connection", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}
