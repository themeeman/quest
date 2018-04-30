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
	"encoding/json"
	commands "discordcommands"
	_ "database/sql"
	command "../Quest/MyCommands"
	"github.com/jmoiron/sqlx"
)

var QuestCommands = commands.HandlerMap{
	"help":        command.Help,
	"mute":        command.Mute,
	"unmute":      command.Unmute,
	"purge":       command.Purge,
	"types":       command.Types,
	"conv":        command.Conv,
	"commit":      command.Commit,
	"addexp":      command.AddExp,
	"setmuterole": command.SetMuteRole,
}
var CommandsData commands.CommandMap
var RegexVerifiers = map[string]string{}

var Token string
var db *sqlx.DB
var guilds []*commands.Guild

const (
	prefix      = "q:"
	commandFile = "Commands.json"
	typesFile   = "Types.json"
)

func ready(s *discordgo.Session, _ *discordgo.Ready) {
	var err error
	guilds, err = commands.QueryAllBotData(s, db, 10000)
	if err != nil {
		log.Println(err)
	}
	for _, v := range guilds {
		fmt.Println(v)
	}
	err = commands.PostAllGuildData(db, guilds)
	if err != nil {
		log.Println(err)
	}
	s.UpdateStatus(0, "q:help")
}

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	fmt.Println(message.Content)
	if strings.HasPrefix(message.Content, prefix) {
		ctx := commands.Bot{
			HandlerMap: &QuestCommands,
			CommandMap: &CommandsData,
			Regex:      RegexVerifiers,
			Prefix:     prefix,
			Guilds:     guilds,
			DB:         db,
			Embed:      questEmbed,
		}
		fmt.Println(guilds)
		err := commands.ExecuteCommand(session, message, ctx)
		if err != nil {
			session.ChannelMessageSendEmbed(message.ChannelID, commands.ErrorEmbed(err))
		}
	}
	fmt.Println("-----------------------")
}

func questEmbed(title string, description string, fields []*discordgo.MessageEmbedField) *discordgo.MessageEmbed {
	emb := &discordgo.MessageEmbed{
		Type:      "rich",
		Title:     title,
		Timestamp: commands.ConvertTimeToTimestamp(time.Now()),
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

func unmarshalJsonFromFile(filename string, v interface{}) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		return
	}
	data := make([]byte, stat.Size())
	_, err = f.Read(data)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, v)
	if err != nil {
		return
	}
	return nil
}

func init() {
	var err error
	err = unmarshalJsonFromFile(commandFile, &CommandsData)
	if err != nil {
		log.Fatalln(err)
	}
	err = unmarshalJsonFromFile(typesFile, &RegexVerifiers)
	if err != nil {
		log.Fatalln(err)
	}
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
	db, err = commands.InitDB()
	if err != nil {
		log.Println(err)
	}
}

func main() {
	defer db.Close()
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
	defer dg.Close()
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

}
