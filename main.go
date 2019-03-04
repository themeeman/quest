package main

import (
	"github.com/tomvanwoow/discordcommands"

	quest "./commands"
	database "./quest/db"
	"./events"
	"./inventory"
	"./structures"
	_ "database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jmoiron/sqlx"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var CommandsData commands.CommandMap
var Types = map[string]string{}

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
var app App
var chests inventory.Chests
var bot *quest.Bot

const (
	prefix = "q:"
)

func questEmbed(title string, description string, fields []*discordgo.MessageEmbedField) *discordgo.MessageEmbed {
	emb := &discordgo.MessageEmbed{
		Type:      "rich",
		Title:     title,
		Timestamp: commands.TimeToTimestamp(time.Now().UTC()),
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

func errorEmbed(e error) *discordgo.MessageEmbed {
	emb := &discordgo.MessageEmbed{
		Type:        "rich",
		Title:       "An error has occurred",
		Description: e.Error(),
		Timestamp:   commands.TimeToTimestamp(time.Now()),
		Color:       0x660000,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Quest Bot",
		},
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
	err = unmarshalJson(app.Types, &Types)
	if err != nil {
		panic(err)
	}
	err = unmarshalJson("src/json/chests.json", &chests)
	if err != nil {
		panic(err)
	}
	db, err = database.InitDB(app.User, app.Pass, app.Host, app.Database)
	if err != nil {
		panic(err)
	}
}

func main() {
	defer db.Close()
	bot = &quest.Bot{
		CommandMap: CommandsData,
		ExpTimes: make(map[struct {
			Guild  string
			Member string
		}]time.Time),
		Errors: make(chan struct {
			Err error
			*discordgo.MessageCreate
		}),
		Types: Types,
		GroupNames: map[commands.Group]string{
			quest.PermissionMember:    "Member",
			quest.PermissionModerator: "Moderator",
			quest.PermissionAdmin:     "Admin",
			quest.PermissionOwner:     "Owner",
		},
		Prefix:     prefix,
		Guilds:     structures.Guilds{},
		DB:         db,
		Chests:     chests,
		Embed:      questEmbed,
		ErrorEmbed: errorEmbed,
	}
	dg, err := commands.NewSession(bot, app.Token)
	if err != nil {
		log.Fatalln("Error making discord session", err)
	}
	e := events.BotEvents{Bot: bot}
	dg.AddHandlerOnce(e.Ready)
	dg.AddHandler(func(session *discordgo.Session, _ *discordgo.Ready) { session.UpdateStatus(0, "q:help") })
	dg.AddHandler(e.MessageCreate)
	dg.AddHandler(e.MemberAdd)
	dg.AddHandler(e.GuildCreate)
	dg.StateEnabled = true
	dg.State.TrackMembers = true
	err = dg.Open()
	if err != nil {
		log.Fatalln("Error opening connection", err)
	}
	defer dg.Close()
	fmt.Println("Quest is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	err = database.PostAllData(db, bot.Guilds)
	if err != nil {
		fmt.Println(err)
		err = database.PostAllData(db, bot.Guilds)
	}
}
