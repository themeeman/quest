package commands

import (
	commands "../../discordcommands"
	"../modlog"
	"fmt"
	"github.com/bwmarrin/discordgo"
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
		return UserNotFoundError{}
	}
	var err error
	if args["Reason"] == "" {
		err = session.GuildMemberDelete(commands.MustGetGuildID(session, message), id)
	} else {
		err = session.GuildMemberDeleteWithReason(commands.MustGetGuildID(session, message), id, args["Reason"])
	}
	if err != nil {
		return fmt.Errorf("Can't kick that user! Make sure I have the discord kick permission.")
	}
	guild := bot.Guilds.Get(commands.MustGetGuildID(session, message))
	if guild.Modlog.Valid {
	}
	session.MessageReactionAdd(message.ChannelID, message.ID, "☑")
	if guild.Modlog.Valid {
		guild.Modlog.Log <- &modlog.CaseKick{
			ModeratorID: message.Author.ID,
			UserID:      id,
			Reason:      args["Reason"],
		}
	}
	return nil
}
