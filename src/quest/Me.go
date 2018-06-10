package MyCommands

import (
	"github.com/bwmarrin/discordgo"
	commands "discordcommands"
	"strconv"
	"fmt"
)

func Me(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string, bot commands.Bot) commands.BotError {
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
	guild := bot.Guilds.Get(commands.MustGetGuildID(session, message))
	m, err := session.GuildMember(guild.ID, id)
	if err != nil {
		return commands.UserNotFoundError{}
	}
	member := guild.Members.Get(id)
	g, _ := session.Guild(commands.MustGetGuildID(session, message))
	rank := commands.GetPermissionLevel(m, guild, g.OwnerID)
	s := []string{"Member", "Moderator", "Admin", "Owner"}
	title := fmt.Sprintf("User %s#%s", m.User.Username, m.User.Discriminator)
	session.ChannelMessageSendEmbed(message.ChannelID, bot.Embed(title, "", []*discordgo.MessageEmbedField{
		{
			Name:  "Experience",
			Value: strconv.Itoa(int(member.Experience)),
		},
		{
			Name:  "Rank",
			Value: s[rank],
		},
	}))
	return nil
}
