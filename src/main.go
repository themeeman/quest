package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"encoding/json"
	commands "./discordcommands"
	_ "database/sql"
	quest "./quest/commands/guild"
	"./quest/commands/direct"
	"./quest/events"
	"github.com/jmoiron/sqlx"
	"math/rand"
	"flag"
)

var QuestCommands = commands.HandlerMap{
	"help":        quest.Help,
	"mute":        quest.Mute,
	"unmute":      quest.Unmute,
	"purge":       quest.Purge,
	"types":       quest.Types,
	"commit":      quest.Commit,
	"addexp":      quest.AddExp,
	"me":          quest.Me,
	"tryparse":    quest.TryParse,
	"massrole":    quest.MassRole,
	"addrole":     quest.AddRole,
	"roles":       quest.Roles,
	"set":         quest.Set,
	"leaderboard": quest.Leaderboard,
	"pull": 	   quest.Pull,
}
var DirectCommands = commands.HandlerMap{
	"help": direct.Help,
}
var CommandsData commands.CommandMap
var RegexVerifiers = map[string]string{}

type App struct {
	Token    string
	User     string
	Pass     string
	Host     string
	Database string
	Commands string
	Types    string
}

var db *sqlx.DB
var guilds commands.Guilds
var app App
var bot *commands.Bot

const (
	prefix = "q:"
)

func questEmbed(title string, description string, fields []*discordgo.MessageEmbedField) *discordgo.MessageEmbed {
	emb := &discordgo.MessageEmbed{
		Type:      "rich",
		Title:     title,
		Timestamp: commands.TimeToTimestamp(time.Now()),
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

func unmarshalJson(filename string, v interface{}) (err error) {
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
	rand.Seed(time.Now().UTC().UnixNano())
	var err error
	if err != nil {
		panic(err)
	}
	var src string
	flag.StringVar(&src, "a", "", "App Location")
	flag.Parse()
	err = unmarshalJson(src, &app)
	if err != nil {
		panic(err)
	}
	err = unmarshalJson(app.Commands, &CommandsData)
	if err != nil {
		panic(err)
	}
	err = unmarshalJson(app.Types, &RegexVerifiers)
	if err != nil {
		panic(err)
	}
	db, err = commands.InitDB(app.User, app.Pass, app.Host, app.Database)
	if err != nil {
		panic(err)
	}
}

func main() {
	defer db.Close()
	bot = &commands.Bot{
		HandlerMap: QuestCommands,
		CommandMap: CommandsData,
		ExpTimes: make(map[struct {
			Guild  string
			Member string
		}]time.Time),
		Regex:  RegexVerifiers,
		Prefix: prefix,
		Guilds: guilds,
		DB:     db,
		Embed:  questEmbed,
	}
	dg, err := discordgo.New("Bot " + app.Token)
	if err != nil {
		log.Fatalln("Error making discord session", err)
		return
	}
	dg.AddHandler(events.Ready(bot))
	dg.AddHandler(events.MessageCreate(bot))
	dg.AddHandler(events.MemberAdd(bot))
	err = dg.Open()
	if err != nil {
		log.Fatalln("Error opening connection", err)
		return
	}
	defer dg.Close()
	fmt.Println("Quest is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	err = commands.PostAllData(db, guilds)
	if err != nil {
		err = commands.PostAllData(db, guilds)
	}
}
