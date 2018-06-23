package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"syscall"
	"strings"
	"time"
	"encoding/json"
	commands "./discordcommands"
	_ "database/sql"
	quest "./quest/commands/guild"
	direct "./quest/commands/direct"
	"github.com/jmoiron/sqlx"
	"math/rand"
	"runtime/debug"
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
var bot commands.Bot

const (
	prefix = "q:"
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
			log.Println(string(debug.Stack()))
			session.ChannelMessageSend(message.ChannelID, "```"+ `An unexpected panic occured in the execution of that command.
`+ fmt.Sprint(r) + "\nTry again later, or contact themeeman#8354" + "```")
		}
	}()
	fmt.Println(message.Content)
	if !message.Author.Bot {
		if strings.ToLower(message.Content) == "good bot" {
			m, _ := session.ChannelMessageSend(message.ChannelID, "Your compliments mean nothing to me")
			time.Sleep(5 * time.Second)
			if m != nil {
				session.ChannelMessageDelete(message.ChannelID, m.ID)
			}
		}
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
	member := g.Members.Get(s.Member)
	if !ok || uint16(time.Since(t).Seconds()) > g.ExpReload {
		bot.ExpTimes[s] = time.Now()
		var r int64
		if g.ExpGainLower > g.ExpGainUpper {
			r = int64(rand.Intn(int(g.ExpGainLower+1-g.ExpGainUpper)) + int(g.ExpGainUpper))
		} else {
			r = int64(rand.Intn(int(g.ExpGainUpper+1-g.ExpGainLower)) + int(g.ExpGainLower))
		}
		member.Experience += r
		fmt.Println(s.Member, r)
		var a int
		if g.LotteryChance == 0 {
			a = 1
		} else {
			a = rand.Intn(int(g.LotteryChance))
		}
		fmt.Println(a)
		if a == 0 {
			ch, err := session.UserChannelCreate(s.Member)
			u, _ := session.GuildMember(s.Guild, s.Member)
			var r int64
			if g.LotteryLower > g.LotteryUpper {
				r = int64(rand.Intn(int(g.LotteryLower+1-g.LotteryUpper)) + int(g.LotteryUpper))
			} else {
				r = int64(rand.Intn(int(g.LotteryUpper+1-g.LotteryLower)) + int(g.LotteryLower))
			}
			if err == nil {
				session.ChannelMessageSend(ch.ID, fmt.Sprintf(`Looks like SOMEBODY is a lucky winner!
That's right, **%s#%s**, you won a grand total of %d Experience! You should give yourself a pat on the back, you're a real winner in life!
🎉🎉🎉🎉🎉🎉🎉`, u.User.Username, u.User.Discriminator, r))
			}
			member.Experience += r
		}
	}

	commands.GrantRoles(session, message, g, member)
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
	dg, err := discordgo.New("Bot " + app.Token)
	if err != nil {
		log.Fatalln("Error making discord session", err)
		return
	}
	dg.AddHandler(ready)
	dg.AddHandler(messageCreate)
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
	err = commands.PostAllData(db, guilds)
	if err != nil {
		err = commands.PostAllData(db, guilds)
	}
}
