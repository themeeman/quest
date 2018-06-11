package discordcommands

import (
	"github.com/bwmarrin/discordgo"
)

func GetPermissionLevel(member *discordgo.Member, guild *Guild, ownerID string) int {
	if member.User.ID == ownerID {
		return 3
	}
	for _, r := range member.Roles {
		if r == guild.AdminRole.String {
			return 2
		} else if r == guild.ModRole.String {
			return 1
		}
	}
	return 0
}

func SufficentPermissions(session *discordgo.Session, message *discordgo.MessageCreate, bot Bot, command *Command) (bool, int, int) {
	g, _ := session.Guild(MustGetGuildID(session, message))
	member, _ := session.GuildMember(g.ID, message.Author.ID)
	guild := bot.Guilds.Get(g.ID)
	had := GetPermissionLevel(member, guild, g.OwnerID)
	required := command.Permission
	return had >= required, had, required
}
