package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/tomvanwoow/quest/utility"
)

func (bot *Bot) RemoveRole(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string) error {
	var roleID string
	if len(message.MentionRoles) > 0 {
		roleID = message.MentionRoles[0]
	} else {
		roleID = args["Role"]
	}
	index := -1
	guild := bot.Guilds.Get(utility.MustGetGuildID(session, message))
	for i, r := range guild.Roles {
		if r.ID == roleID {
			index = i
		}
	}
	if index == -1 {
		return RoleError{roleID}
	}
	guild.Roles = append(guild.Roles[:index], guild.Roles[index+1:]...)
	_ = session.MessageReactionAdd(message.ChannelID, message.ID, "☑")
	return nil
}
