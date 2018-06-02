package MyCommands

import (
	"github.com/bwmarrin/discordgo"
	commands "discordcommands"
	"strconv"
	"fmt"
)

func Me(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string, bot commands.Bot) commands.BotError {
	var id string
	if args["User"] != "" {
		id = message.Mentions[0].ID
	} else {
		id = message.Author.ID
	}
	guild, _ := commands.FindGuildByID(bot.Guilds, commands.MustGetGuildID(session, message))
	member, ok := commands.FindMemberByID(guild.Members, id)
	fmt.Println(member, ok)
	if !ok {
		guild.Members = append(guild.Members, member)
	}
	fmt.Println(guild.Members)
	session.ChannelMessageSendEmbed(message.ChannelID, bot.Embed("User", "", []*discordgo.MessageEmbedField{
		{
			Name:  "Experience",
			Value: strconv.Itoa(member.Experience),
		},
	}))
	return nil
}
