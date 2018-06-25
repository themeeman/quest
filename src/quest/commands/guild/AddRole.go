package guild

import (
	"github.com/bwmarrin/discordgo"
	commands "../../../discordcommands"
	"strconv"
	"fmt"
)

func AddRole(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string, bot *commands.Bot) commands.BotError {
	var roleID string
	if len(message.MentionRoles) > 0 {
		roleID = message.MentionRoles[0]
	} else {
		roleID = args["Role"]
	}
	guild := bot.Guilds.Get(commands.MustGetGuildID(session, message))
	exp, _ := strconv.Atoi(args["Experience"])
	role := &commands.Role{
		Experience: int64(exp),
		ID:         roleID,
	}
	allIDs := make([]string, len(guild.Roles))
	for i, v := range guild.Roles {
		allIDs[i] = v.ID
	}
	ok, index := commands.Contains(allIDs, roleID)
	fmt.Println(ok, index)
	if ok {
		guild.Roles[index] = role
	} else if len(guild.Roles) >= 64 {
		return commands.CustomError("Invalid action - 64 roles is the absolute limit\nTry removing a role")
	} else {
		guild.Roles = append(guild.Roles, role)
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
