package guild

import (
	"github.com/bwmarrin/discordgo"
	commands "../../../discordcommands"
	"strconv"
)

func AddRole(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string, bot commands.Bot) commands.BotError {
	var role string
	if len(message.MentionRoles) > 0 {
		role = message.MentionRoles[0]
	} else {
		role = args["Role"]
	}
	guild := bot.Guilds.Get(commands.MustGetGuildID(session, message))
	if isRoleIn(guild.Roles, role) {
		return commands.UnknownError{}
	}
	exp, _ := strconv.Atoi(args["Experience"])
	guild.Roles = append(guild.Roles, &commands.Role{
		Experience: int64(exp),
		ID:         role,
	})
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
