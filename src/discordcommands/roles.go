package discordcommands

import (
	"github.com/bwmarrin/discordgo"
	"errors"
)

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
