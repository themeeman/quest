package MyCommands

import (
	"github.com/bwmarrin/discordgo"
	commands "discordcommands"
	"strconv"
	"fmt"
)

func AddExp(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string, bot commands.Bot) commands.BotError {
	c, _ := session.Channel(message.ChannelID)
	newGuilds, guild, index, ok := commands.FindGuildByID(bot.Guilds, c.GuildID)
	fmt.Println(newGuilds, guild, index, ok)
	var id string
	if args["User"] == "" {
		id = message.Author.ID
	} else {
		id = message.Mentions[0].ID
	}
	newMembers, member, index, ok := commands.FindMemberByID(guild.Members, id)
	fmt.Println(newMembers, member, index, ok)
	guild.Members = newMembers
	exp, _ := strconv.Atoi(args["Value"])
	member.Experience += exp
	session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Total Experience: %d", guild.Members[index].Experience))
	return nil
}
