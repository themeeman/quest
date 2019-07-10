package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/tomvanwoow/quest/modlog"
	"github.com/tomvanwoow/quest/utility"
)

func (bot *Bot) Kick(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string) error {
	var id string
	if args["User"] == "" {
		id = message.Author.ID
	} else if len(args["User"]) == 18 {
		id = args["User"]
	} else if len(message.Mentions) > 0 {
		id = message.Mentions[0].ID
	} else {
		return UserNotFoundError
	}
	var err error
	if args["Reason"] == "" {
		err = session.GuildMemberDelete(utility.MustGetGuildID(session, message), id)
	} else {
		err = session.GuildMemberDeleteWithReason(utility.MustGetGuildID(session, message), id, args["Reason"])
	}
	if err != nil {
		return fmt.Errorf("Can't kick that user! Make sure I have the discord kick permission.")
	}
	guild := bot.Guilds.Get(utility.MustGetGuildID(session, message))
	session.MessageReactionAdd(message.ChannelID, message.ID, "☑")
	guild.RLock()
	if guild.Modlog.Valid {
		guild.Modlog.Log <- &modlog.CaseKick{
			ModeratorID: message.Author.ID,
			UserID:      id,
			Reason:      args["Reason"],
		}
	}
	guild.RUnlock()
	return nil
}
