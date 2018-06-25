package discordcommands

import (
	"github.com/bwmarrin/discordgo"
)

func GetPermissionLevel(session *discordgo.Session, member *discordgo.Member, guild *Guild, ownerID string) int {
	if member.User.ID == ownerID {
		return 3
	}
	for _, r := range member.Roles {
		role, err := FindRole(session, guild.ID, r)
		if err == nil {
			if role.Permissions & discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {
				return 2
			}
		}
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

func SufficentPermissions(session *discordgo.Session, guild *discordgo.Guild, member *discordgo.Member, bot *Bot, command *Command) (bool, int, int) {
	g := bot.Guilds.Get(guild.ID)
	had := GetPermissionLevel(session, member, g, guild.OwnerID)
	required := command.Permission
	return had >= required, had, required
}
