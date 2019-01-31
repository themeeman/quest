package commands

import (
	commands "../../discordcommands"
	"../structures"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"../utility"
)

func (bot *Bot) AddRole(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string) error {
	var roleID string
	if len(message.MentionRoles) > 0 {
		roleID = message.MentionRoles[0]
	} else {
		roleID = args["Role"]
	}
	guild := bot.Guilds.Get(commands.MustGetGuildID(session, message))
	exp, _ := strconv.Atoi(args["Experience"])
	if exp == 0 {
		return fmt.Errorf(`I see you are trying to add a role for 0 experience!
If you want to, it is better to use q:set Autorole <role>`)
	}
	role := &structures.Role{
		Experience: int64(exp),
		ID:         roleID,
	}
	allIDs := make([]string, len(guild.Roles))
	for i, v := range guild.Roles {
		allIDs[i] = v.ID
	}
	ok, index := utility.Contains(allIDs, roleID)
	fmt.Println(ok, index)
	if ok {
		guild.Roles[index] = role
	} else if len(guild.Roles) >= 64 {
		return fmt.Errorf("Invalid action - 64 roles is the absolute limit\nTry removing a role")
	} else {
		guild.Roles = append(guild.Roles, role)
	}
	session.MessageReactionAdd(message.ChannelID, message.ID, "☑")
	return nil
}