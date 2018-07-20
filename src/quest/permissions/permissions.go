package permissions

import (
	commands "../../discordcommands"
	"github.com/bwmarrin/discordgo"
	"../structures"
)

const (
	PermissionMember    commands.Group = iota
	PermissionModerator
	PermissionAdmin
	PermissionOwner
)

func GetPermissionLevel(session *discordgo.Session, member *discordgo.Member, guild *structures.Guild, ownerID string) commands.Group {
	if member.User.ID == ownerID {
		return PermissionOwner
	}
	for _, r := range member.Roles {
		role, err := commands.GetRole(session, guild.ID, r)
		if err == nil {
			if role.Permissions & discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {
				return PermissionAdmin
			}
		}
	}
	for _, r := range member.Roles {
		if r == guild.AdminRole.String {
			return PermissionAdmin
		} else if r == guild.ModRole.String {
			return PermissionModerator
		}
	}
	return PermissionMember
}
