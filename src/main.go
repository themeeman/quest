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
	quest "./quest"
	"github.com/jmoiron/sqlx"
	"math/rand"
	"runtime/debug"
)

var QuestCommands = commands.HandlerMap{
	"help":     quest.Help,
	"mute":     quest.Mute,
	"unmute":   quest.Unmute,
	"purge":    quest.Purge,
	"types":    quest.Types,
	"commit":   quest.Commit,
	"addexp":   quest.AddExp,
	"me":       quest.Me,
	"tryparse": quest.TryParse,
	"massrole": quest.MassRole,
	"addrole":  quest.AddRole,
	"roles":    quest.Roles,
	"set":      quest.Set,
}
var CommandsData commands.CommandMap
var RegexVerifiers = map[string]string{}

var Token string
var db *sqlx.DB
var guilds commands.Guilds

var bot commands.Bot

const (
	prefix      = "q:"
	commandFile = `src\json\Commands.json`
	typesFile   = `src\json\Types.json`
)

func ready(s *discordgo.Session, _ *discordgo.Ready) {
	var err error
	guilds, err = commands.QueryAllData(db)
	if err != nil {
		log.Println("b", err)
	}
	for _, v := range guilds {
		fmt.Println(v)
	}
	s.UpdateStatus(0, "q:help")
	bot = commands.Bot{
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
	go func() {
		for {
			time.Sleep(time.Minute * 10)
			err := commands.PostAllData(db, bot.Guilds)
			if err != nil {
				log.Println(err)
			} else {
				log.Println("Successfully commited all data")
			}
		}
	}()
}

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(debug.Stack())
			session.ChannelMessageSend(message.ChannelID, "```"+ `An unexpected panic occured in the execution of that quest.
`+ fmt.Sprint(r)+ "```")
		}
	}()
	fmt.Println(message.Content)
	if !message.Author.Bot {
		if strings.HasPrefix(strings.ToLower(message.Content), prefix) {
			err := commands.ExecuteCommand(session, message, bot)
			if err != nil {
				session.ChannelMessageSendEmbed(message.ChannelID, commands.ErrorEmbed(err))
			}
		}
		grantExp(&bot, session, message)
		fmt.Println("-----------------------")
	}
}

func guildCreate(_ *discordgo.Session, event *discordgo.GuildCreate) {
	commands.CreateAllData(db, event.Guild.ID)
}

func memberAdd(session *discordgo.Session, event *discordgo.GuildMemberAdd) {
	guild := bot.Guilds.Get(event.GuildID)
	if guild.Autorole.Valid {
		session.GuildMemberRoleAdd(event.GuildID, event.Member.User.ID, guild.Autorole.String)
	}
}

func grantExp(bot *commands.Bot, session *discordgo.Session, message *discordgo.MessageCreate) {
	s := struct {
		Guild  string
		Member string
	}{
		Guild:  commands.MustGetGuildID(session, message),
		Member: message.Author.ID,
	}
	t, ok := bot.ExpTimes[s]
	g := bot.Guilds.Get(s.Guild)
	m := g.Members.Get(s.Member)
	fmt.Println(s.Member)
	if !ok || uint16(time.Since(t).Seconds()) > g.ExpReload {
		bot.ExpTimes[s] = time.Now()
		m.Experience += int64(rand.Intn(10) + 10)
	}
	commands.GrantRoles(session, message, g, m)
}

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
	err = unmarshalJson(commandFile, &CommandsData)
	if err != nil {
		panic(err)
	}
	err = unmarshalJson(typesFile, &RegexVerifiers)
	if err != nil {
		panic(err)
	}
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
	db, err = commands.InitDB()
	if err != nil {
		panic(err)
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
	dg.AddHandler(guildCreate)
	dg.AddHandler(memberAdd)
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
	commands.PostAllData(db, guilds)
}
