package experience

import (
	"github.com/bwmarrin/discordgo"
	"time"
	commands "../../discordcommands"
	"fmt"
)

func GrantRoles(session *discordgo.Session, message *discordgo.MessageCreate, guild *commands.Guild, member *commands.Member) error {
	m, _ := session.GuildMember(commands.MustGetGuildID(session, message), member.ID)
	for _, r := range guild.Roles {
		if m != nil && member.Experience >= r.Experience {
			role, err := commands.FindRole(session, guild.ID, r.ID)
			if err != nil {
				continue
			}
			found, _ := commands.Contains(m.Roles, role.ID)
			if !found {
				session.GuildMemberRoleAdd(guild.ID, member.ID, r.ID)
				session.ChannelMessageSendEmbed(message.ChannelID, questEmbedColor(role.Name,
					fmt.Sprintf("%s#%s Received Role", m.User.Username, m.User.Discriminator), nil, role.Color))
			}
		}
	}
	return nil
}

func questEmbedColor(title string, description string, fields []*discordgo.MessageEmbedField, color int) *discordgo.MessageEmbed {
	emb := &discordgo.MessageEmbed{
		Type:      "rich",
		Title:     title,
		Timestamp: commands.TimeToTimestamp(time.Now()),
		Color:     color,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Quest Bot",
		},
		Fields: fields,
	}
	if description != "" {
		emb.Description = description
	}
	return emb
}
