package guild

import (
	"github.com/bwmarrin/discordgo"
	commands "../../../discordcommands"
	"strconv"
)

func AddRole(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string, bot *commands.Bot) commands.BotError {
	var role string
	if len(message.MentionRoles) > 0 {
		role = message.MentionRoles[0]
	} else {
		role = args["Role"]
	}
	guild := bot.Guilds.Get(commands.MustGetGuildID(session, message))
	exp, _ := strconv.Atoi(args["Experience"])
	ok, index := commands.Contains(guild.Roles, role)
	if ok {
		guild.Roles[index] = &commands.Role{
			Experience: int64(exp),
			ID:         role,
		}
	} else if len(guild.Roles) >= 64 {
		return commands.CustomError("Invalid action - 64 roles is the absolute limit\nTry removing a role")
	} else {
		guild.Roles = append(guild.Roles, &commands.Role{
			Experience: int64(exp),
			ID:         role,
		})
	}
	session.MessageReactionAdd(message.ChannelID, message.ID, "☑")
	return nil
}

func isRoleIn(roles commands.Roles, id string) bool {
	for _, r := range roles {
		if r.ID == id {
			return true
		}
	}
	return false
}
