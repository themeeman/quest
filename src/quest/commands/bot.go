package commands

import (
	commands "../../discordcommands"
	"../inventory"
	"../structures"
	"github.com/bwmarrin/discordgo"
	"github.com/jmoiron/sqlx"
	"time"
)

type Bot struct {
	commands.CommandMap
	ExpTimes map[struct {
		Guild  string
		Member string
	}]time.Time
	Types      map[string]string
	GroupNames map[commands.Group]string
	Prefix     string
	structures.Guilds
	DB     *sqlx.DB
	Errors chan struct {
		Err error
		*discordgo.MessageCreate
	}
	inventory.Chests
	ErrorEmbed func(e error) *discordgo.MessageEmbed
	Embed func(title string,
		description string,
		fields []*discordgo.MessageEmbedField) *discordgo.MessageEmbed
}
