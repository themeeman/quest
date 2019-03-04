package commands

import (
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/tomvanwoow/quest/structures"
	"github.com/tomvanwoow/quest/utility"
	"sort"
	"time"
)

func (bot *Bot) Roles(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string) error {
	guild := bot.Guilds.Get(utility.MustGetGuildID(session, message))
	if len(guild.Roles) == 0 {
		session.ChannelMessageSend(message.ChannelID, "No reward roles configured\nUse q:addrole to create some")
		return nil
	}
	roles := make(structures.Roles, len(guild.Roles))
	copy(roles, guild.Roles)
	sort.Sort(roles)
	var buffer bytes.Buffer
	buffer.WriteString("```\n")
	allRoles, err := session.GuildRoles(guild.ID)
	if err != nil {
		return nil
	}
	discordRoles := make(discordgo.Roles, len(roles))
	for i, r := range roles {
		ok, index := roleContains(allRoles, r.ID)
		if ok {
			discordRoles[i] = allRoles[index]
		}
	}
	for i, v := range discordRoles {
		r := roles[i]
		if r != nil && v != nil {
			buffer.WriteString(fmt.Sprintf("%v EXP: %s\n", r.Experience, v.Name))
		}
	}
	buffer.WriteString("```")
	t := time.Now()
	session.ChannelMessageSend(message.ChannelID, buffer.String())
	fmt.Println(time.Since(t))
	return nil
}

func roleContains(roles []*discordgo.Role, id string) (bool, int) {
	for i, r := range roles {
		if r.ID == id {
			return true, i
		}
	}
	return false, 0
}
