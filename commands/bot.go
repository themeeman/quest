package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jmoiron/sqlx"
	commands "github.com/tomvanwoow/disgone"
	"github.com/tomvanwoow/quest/inventory"
	"github.com/tomvanwoow/quest/structures"
	"time"
)

type Bot struct {
	*commands.BotOptions
	ExpTimes map[struct {
		Guild  string
		Member string
	}]time.Time
	structures.Guilds
	DB *sqlx.DB
	inventory.Chests
	ErrorEmbed func(e error) *discordgo.MessageEmbed
	Embed      func(title string,
		description string,
		fields []*discordgo.MessageEmbedField) *discordgo.MessageEmbed
}
