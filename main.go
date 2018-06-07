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
	"math/rand"
	"runtime/debug"
)

var QuestCommands = commands.HandlerMap{
	"help":     command.Help,
	"mute":     command.Mute,
	"unmute":   command.Unmute,
	"purge":    command.Purge,
	"types":    command.Types,
	"commit":   command.Commit,
	"addexp":   command.AddExp,
	"setmute":  command.SetMuteRole,
	"me":       command.Me,
	"tryparse": command.TryParse,
	"massrole": command.MassRole,
	"addrole":  command.AddRole,
	"roles":    command.Roles,
}
var CommandsData commands.CommandMap
var RegexVerifiers = map[string]string{}

var Token string
var db *sqlx.DB
var guilds commands.Guilds

var ctx commands.Bot

const (
	prefix      = "q:"
	commandFile = "Commands.json"
	typesFile   = "Types.json"
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
	ctx = commands.Bot{
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
}

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	defer func() {
		if r := recover(); r != nil {
			rep := strings.Replace
			text := fmt.Sprint(r) + "\n" + rep(rep(string(debug.Stack()), "Tom Van Wowe", "Root", -1), "Tom_Van_Wowe", "Root", -1)
			session.ChannelMessageSend(message.ChannelID, "```"+ `An unexpected panic occured in the execution of that command.
`+ text+ "```")
		}
	}()
	fmt.Println(message.Content)
	if !message.Author.Bot {
		if strings.HasPrefix(strings.ToLower(message.Content), prefix) {
			err := commands.ExecuteCommand(session, message, ctx)
			if err != nil {
				session.ChannelMessageSendEmbed(message.ChannelID, commands.ErrorEmbed(err))
			}
		}
		grantExp(&ctx, session, message)
		fmt.Println("-----------------------")
	}
}

func guildCreate(_ *discordgo.Session, event *discordgo.GuildCreate) {
	commands.CreateAllData(db, event.Guild.ID)
}

func memberAdd(session *discordgo.Session, event *discordgo.GuildMemberAdd) {
	guild := ctx.Guilds.Get(event.GuildID)
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
	g := ctx.Guilds.Get(s.Guild)
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
