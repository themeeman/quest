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
)

type ArgumentData struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Optional bool   `json:"optional"`
	Infinite bool   `json:"infinite"`
}

type CommandData struct {
	Description   string         `json:"description"`
	ArgumentsData []ArgumentData `json:"arguments"`
	Cooldown      int            `json:"cooldown"`
	Permission    int            `json:"permission"`
}

type Command func(session *discordgo.Session, message *discordgo.MessageCreate, args []string)

type Commands map[string]Command

var QuestCommands = Commands{
	"help": Help,
	"mute": Mute,
}
var CommandsData map[string]CommandData
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

func (cd CommandData) CalcForcedArgs() int {
	var i int
	for _, v := range cd.ArgumentsData {
		if !v.Optional {
			i += 1
		}
	}
	return i
}
func ready(s *discordgo.Session, _ *discordgo.Ready) {
	s.UpdateStatus(0, "q:")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println(m.Content)
	if strings.HasPrefix(m.Content, prefix) {
		cmd := strings.TrimPrefix(m.Content, prefix)
		args := strings.Split(cmd, " ")
		ExecuteCommand(s, m, args)
	}
}

func ExecuteCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	for i, v := range args {
		fmt.Println(i, v)
	}
	cmdName := args[0]
	cmd, ok := QuestCommands[cmdName]
	if !ok {
		s.ChannelMessageSend(m.ChannelID, "Invalid command")
		return
	}
	cmdInfo, ok := CommandsData[cmdName]
	if !ok {
		return
	}
	if cmdInfo.CalcForcedArgs() > 0 && len(args) == 1 {
		Help(s, m, []string{"help", cmdName})
		return
	}
	cmd(s, m, args)

}

func Help(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) == 1 {
		var value string
		for name, v := range CommandsData {
			value += fmt.Sprintf("**%v - ** %v\n", name, v.Description)
		}
		fields := []*discordgo.MessageEmbedField{
			{
				"Commands",
				value,
				false,
			},
		}
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, QuestEmbed("Help", "", fields))
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		cmdName := args[1]
		cmdInfo, ok := CommandsData[cmdName]
		if !ok {
			sss := fmt.Sprintf("The command %s does not exist.\nUse the help command for more info.", cmdName)
			_, err := s.ChannelMessageSendEmbed(m.ChannelID, QuestEmbed("Command Not Found!", sss, nil))
			if err != nil {
				log.Println(err)
			}
			return
		}
		var usage string
		usage = prefix + cmdName + " "
		for _, v := range cmdInfo.ArgumentsData {
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
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, QuestEmbed(strings.ToTitle(cmdName), `\`+cmdInfo.Description, fields))
		if err != nil {
			log.Println(err)
		}
	}
}

func Mute(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) >= 3 {
		ch, err := s.State.Channel(m.ChannelID)
		if err != nil {
			log.Println("Could not find channel to mute")
			return
		}
		guild, err := s.Guild(ch.GuildID)
		if err != nil {
			log.Println("Could not find guild")
			return
		}
		dur, err := strconv.Atoi(args[2])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Expected integer, got something else")
			return
		}
		if len(m.Mentions) != 1 {
			s.ChannelMessageSend(m.ChannelID, "Nope")
			return
		}
		user := m.Mentions[0]
		for _, u := range guild.Members {
			for _, r := range u.Roles {
				if r == "413273250131345409" {
					s.ChannelMessageSend(m.ChannelID, "Could not mute user, they are already muted!")
					return
				}
			}
		}
		err = s.GuildMemberRoleAdd(ch.GuildID, user.ID, "413273250131345409")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "User not found or something")
			return
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Successfully muted %s!", user.Username))
		go func() {
			time.Sleep(time.Second * time.Duration(dur))
			s.GuildMemberRoleRemove(ch.GuildID, m.Mentions[0].ID, "413273250131345409")
		}()
	}
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

func QuestErrorEmbed() {}

func ConvertTimeToTimestamp(t time.Time) string {
	return t.Format("2006-01-02T15:04:05+00:00")
}

func init() {
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
