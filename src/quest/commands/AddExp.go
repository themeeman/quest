package commands

import (
	commands "../../discordcommands"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
)

func (bot *Bot) AddExp(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string) error {
	guild := bot.Guilds.Get(commands.MustGetGuildID(session, message))
	var id string
	if args["User"] == "" {
		id = message.Author.ID
	} else if len(args["User"]) == 18 {
		id = args["User"]
	} else if len(message.Mentions) > 0 {
		id = message.Mentions[0].ID
	} else {
		return UserNotFoundError{}
	}
	_, err := session.GuildMember(commands.MustGetGuildID(session, message), id)
	if err != nil {
		return UserNotFoundError{}
	}
	member := guild.Members.Get(id)
	exp, _ := strconv.Atoi(strings.Replace(args["Value"], ",", "", -1))
	member.Experience += int64(exp)
	session.MessageReactionAdd(message.ChannelID, message.ID, "☑")
	return nil
}
