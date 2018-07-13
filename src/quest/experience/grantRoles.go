package experience

import (
	"github.com/bwmarrin/discordgo"
	"time"
	commands "../../discordcommands"
	"fmt"
)

func GrantRoles(session *discordgo.Session, message *discordgo.MessageCreate, guild *commands.Guild, member *commands.Member) error {
	m, err := session.GuildMember(guild.ID, member.ID)
	if err != nil {
		return err
	}
	for _, r := range guild.Roles {
		if member.Experience >= r.Experience {
			role, err := commands.FindRole(session, guild.ID, r.ID)
			if err != nil {
				continue
			}
			found, _ := commands.Contains(m.Roles, role.ID)
			if !found {
				session.GuildMemberRoleAdd(guild.ID, member.ID, r.ID)
				session.ChannelMessageSendEmbed(message.ChannelID, questEmbedColor(m.User.Username, m.User.Discriminator, role.Name, role.Color))
			}
		}
	}
	return nil
}

func questEmbedColor(username string, discriminator string, rolename string, color int) *discordgo.MessageEmbed {
	emb := &discordgo.MessageEmbed{
		Type:        "rich",
		Title:       fmt.Sprintf("Congratulations %s#%s", username, discriminator),
		Description: fmt.Sprintf("You received the %s role", rolename),
		Timestamp:   commands.TimeToTimestamp(time.Now()),
		Color:       color,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Quest Bot",
		},
	}
	return emb
}
