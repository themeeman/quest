package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/tomvanwoow/quest/inventory"
	"github.com/tomvanwoow/quest/structures"
	"github.com/tomvanwoow/quest/utility"
	"time"
)

func (bot *Bot) Daily(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string) error {
	guild := bot.Guilds.Get(utility.MustGetGuildID(session, message))
	guild.RLock()
	member := guild.Members.Get(message.Author.ID)
	guild.RUnlock()
	member.RLock()
	defer member.RUnlock()
	if !member.LastDaily.Valid || time.Since(member.LastDaily.Time.UTC()) > 23*time.Hour {
		guild.Members.Apply(member.ID, func(member *structures.Member) {
			member.LastDaily.Valid = true
			member.LastDaily.Time = time.Now().UTC()
			member.Chests[inventory.ChestDaily] += 1
		})
		_, _ = session.ChannelMessageSendEmbed(message.ChannelID, bot.Embed("", "", []*discordgo.MessageEmbedField{
			{
				Name:  "Collected Daily Reward!",
				Value: "+1 **Daily Chest**",
			},
		}))
	} else {
		t := member.LastDaily.Time.Add(23 * time.Hour).Sub(time.Now().UTC())
		_, _ = session.ChannelMessageSend(message.ChannelID,
			fmt.Sprintf("Sorry, you need to wait %d hours and %d minutes before your next redeem",
				int(t.Hours()), int(t.Minutes())%60))
	}
	return nil
}
