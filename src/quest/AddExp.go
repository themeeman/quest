package MyCommands

import (
	"github.com/bwmarrin/discordgo"
	commands ".././discordcommands"
	"strconv"
	"strings"
)

func AddExp(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string, bot commands.Bot) commands.BotError {
	guild := bot.Guilds.Get(commands.MustGetGuildID(session, message))
	var id string
	if args["User"] == "" {
		id = message.Author.ID
	} else if len(args["User"]) == 18 {
		id = args["User"]
	} else if len(message.Mentions) > 0 {
		id = message.Mentions[0].ID
	} else {
		return commands.UserNotFoundError{}
	}
	_, err := session.GuildMember(commands.MustGetGuildID(session, message), id)
	if err != nil {
		return commands.UserNotFoundError{}
	}
	member := guild.Members.Get(id)
	exp, _ := strconv.Atoi(strings.Replace(args["Value"], ",", "", -1))
	member.Experience += int64(exp)
	commands.GrantRoles(session, message, guild, member)
	session.MessageReactionAdd(message.ChannelID, message.ID, "☑")
	return nil
}
