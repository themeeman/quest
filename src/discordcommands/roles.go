package discordcommands

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
	"errors"
	"time"
)

func GrantRoles(session *discordgo.Session, message *discordgo.MessageCreate, guild *Guild, member *Member) error {
	m, _ := session.GuildMember(MustGetGuildID(session, message), member.ID)
	for _, r := range guild.Roles {
		if m != nil && member.Experience >= r.Experience {
			role, err := FindRole(session, guild.ID, r.ID)
			if err != nil {
				continue
			}
			var found bool
			for _, rr := range m.Roles {
				if rr == r.ID {
					found = true
				}
			}
			if !found {
				session.GuildMemberRoleAdd(guild.ID, member.ID, r.ID)
				session.ChannelMessageSendEmbed(message.ChannelID, questEmbedColor(role.Name,
					fmt.Sprintf("%s#%s Received Role", m.User.Username, m.User.Discriminator), nil, role.Color))
			}
		}
	}
	return nil
}

func FindRole(session *discordgo.Session, guildid string, id string) (*discordgo.Role, error) {
	rs, err := session.GuildRoles(guildid)
	if err != nil {
		return nil, err
	}
	for _, r := range rs {
		if r.ID == id {
			return r, nil
		}
	}
	return nil, errors.New("role not found")
}

func questEmbedColor(title string, description string, fields []*discordgo.MessageEmbedField, color int) *discordgo.MessageEmbed {
	emb := &discordgo.MessageEmbed{
		Type:      "rich",
		Title:     title,
		Timestamp: TimeToTimestamp(time.Now()),
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
