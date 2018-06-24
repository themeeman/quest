package guild

import (
	"github.com/bwmarrin/discordgo"
	commands "../../../discordcommands"
	"fmt"
)

func MassRole(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string, bot *commands.Bot) commands.BotError {
	c, _ := session.Channel(message.ChannelID)
	guild := bot.Guilds.Get(c.GuildID)
	var role string
	if args["Role"] == "" {
		if !guild.Autorole.Valid {
			return commands.AutoRoleError{}
		} else {
			role = guild.Autorole.String
		}
	} else if len(args["Role"]) > 18 {
		role = message.MentionRoles[0]
	} else {
		role = args["Role"]
	}
	go func() {
		g, _ := session.Guild(c.GuildID)
		for _, m := range g.Members {
			err := session.GuildMemberRoleAdd(g.ID, m.User.ID, role)
			if err != nil {
				fmt.Println(err)
			}
		}
		session.MessageReactionAdd(message.ChannelID, message.ID, "☑")
	}()
	return nil
}
