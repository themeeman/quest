package commands

import (
	commands "../../discordcommands"
	"../inventory"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

func (bot *Bot) Daily(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string) error {
	guild := bot.Guilds.Get(commands.MustGetGuildID(session, message))
	member := guild.Members.Get(message.Author.ID)
	if !member.LastDaily.Valid || time.Since(member.LastDaily.Time.UTC()) > 23*time.Hour {
		member.LastDaily.Valid = true
		member.LastDaily.Time = time.Now().UTC()
		member.Chests[inventory.ChestDaily] += 1
		session.ChannelMessageSendEmbed(message.ChannelID, bot.Embed("", "", []*discordgo.MessageEmbedField{
			{
				Name:  "Collected Daily Reward!",
				Value: "+1 **Daily Chest**",
			},
		}))
	} else {
		t := member.LastDaily.Time.Add(23 * time.Hour).Sub(time.Now().UTC())
		session.ChannelMessageSend(message.ChannelID,
			fmt.Sprintf("Sorry, you need to wait %d hours and %d minutes before your next redeem",
				int(t.Hours()), int(t.Minutes())%60))
	}
	return nil
}
