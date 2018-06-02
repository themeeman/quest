package MyCommands

import (
	"github.com/bwmarrin/discordgo"
	commands "discordcommands"
	"strconv"
	"fmt"
)

func AddExp(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string, bot commands.Bot) commands.BotError {
	c, _ := session.Channel(message.ChannelID)
	guild, ok := commands.FindGuildByID(bot.Guilds, c.GuildID)
	fmt.Println(guild, ok)
	if !ok {
		bot.Guilds = append(bot.Guilds, guild)
		guild = bot.Guilds[len(bot.Guilds) - 1]
	}
	var id string
	if args["User"] == "" {
		id = message.Author.ID
	} else {
		id = message.Mentions[0].ID
	}
	member, ok := commands.FindMemberByID(guild.Members, id)
	fmt.Println(member, ok)
	if !ok {
		guild.Members = append(guild.Members, member)
	}
	exp, _ := strconv.Atoi(args["Value"])
	member.Experience += exp
	for _, v := range bot.Guilds {
		fmt.Println(v.ID)
		for _, m := range v.Members {
			fmt.Println(*m)
		}
	}
	session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Total Experience: %d", member.Experience))
	return nil
}
