package discordcommands

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
)

type Group uint

func GetRole(session *discordgo.Session, guildID string, id string) (*discordgo.Role, error) {
	rs, err := session.GuildRoles(guildID)
	if err != nil {
		return nil, err
	}
	for _, r := range rs {
		if r.ID == id {
			return r, nil
		}
	}
	return nil, fmt.Errorf("role not found %s", id)
}