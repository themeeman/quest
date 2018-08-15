package commands

import (
	"time"
	"github.com/jmoiron/sqlx"
	"github.com/bwmarrin/discordgo"
	commands "../../discordcommands"
	"../structures"
	"../inventory"
)

type Bot struct {
	commands.CommandMap
	ExpTimes map[struct {
		Guild  string
		Member string
	}]time.Time
	Regex      map[string]string
	GroupNames map[commands.Group]string
	Prefix     string
	structures.Guilds
	DB         *sqlx.DB
	Errors chan struct {
		Err error
		*discordgo.MessageCreate
	}
	inventory.Chests
	Embed func(title string,
		description string,
		fields []*discordgo.MessageEmbedField) *discordgo.MessageEmbed
}
