package guild

import (
	"github.com/bwmarrin/discordgo"
	commands "../../../discordcommands"
	"sort"
	"bytes"
	"fmt"
)

func Roles(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string, bot *commands.Bot) commands.BotError {
	guild := bot.Guilds.Get(commands.MustGetGuildID(session, message))
	if len(guild.Roles) == 0 {
		session.ChannelMessageSend(message.ChannelID, "No reward roles configured\nUse q:addrole to create some")
		return nil
	}
	roles := make(commands.Roles, len(guild.Roles))
	copy(roles, guild.Roles)
	sort.Sort(roles)
	var buffer bytes.Buffer
	buffer.WriteString("```\n")
	for _, v := range roles {
		r, err := commands.FindRole(session, guild.ID, v.ID)
		if err != nil {
			fmt.Println(err)
			continue
		}
		buffer.WriteString(fmt.Sprintf("%v EXP: %s\n", v.Experience, r.Name))
	}
	buffer.WriteString("```")
	session.ChannelMessageSend(message.ChannelID, buffer.String())
	return nil
}
