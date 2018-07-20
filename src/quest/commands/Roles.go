package commands

import (
	"github.com/bwmarrin/discordgo"
	commands "../../discordcommands"
	"../structures"
	"sort"
	"bytes"
	"fmt"
)

func (bot *Bot) Roles(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string) error {
	guild := bot.Guilds.Get(commands.MustGetGuildID(session, message))
	if len(guild.Roles) == 0 {
		session.ChannelMessageSend(message.ChannelID, "No reward roles configured\nUse q:addrole to create some")
		return nil
	}
	roles := make(structures.Roles, len(guild.Roles))
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
