package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/tomvanwoow/quest/modlog"
	"github.com/tomvanwoow/quest/utility"
)

func (bot *Bot) TestLog(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string) error {
	const myID = "164759167561629696"
	if message.Author.ID != myID {
		return nil
	}
	guild := bot.Guilds.Get(utility.MustGetGuildID(session, message))
	guild.Modlog.Log <- &modlog.CaseAddExp{
		ModeratorID: message.Author.ID,
		Experience:  100000,
		UserID:      myID,
	}
	guild.Modlog.Log <- &modlog.CaseBan{
		ModeratorID: message.Author.ID,
		UserID:      myID,
		Reason:      args["Reason"],
	}
	guild.Modlog.Log <- &modlog.CaseKick{
		ModeratorID: message.Author.ID,
		UserID:      myID,
		Reason:      "test",
	}
	guild.Modlog.Log <- &modlog.CaseMute{
		ModeratorID: message.Author.ID,
		UserID:      myID,
		Duration:    100,
		Reason:      "test",
	}
	guild.Modlog.Log <- &modlog.CasePurge{
		ModeratorID: message.Author.ID,
		ChannelID:   message.ChannelID,
		Amount:      100,
	}
	guild.Modlog.Log <- &modlog.CaseSet{
		ModeratorID: message.Author.ID,
		Option:      "test",
		Value:       "test",
	}
	guild.Modlog.Log <- &modlog.CaseUnban{
		ModeratorID: message.Author.ID,
		UserID:      myID,
		Reason:      "test",
	}
	guild.Modlog.Log <- &modlog.CaseUnmute{
		ModeratorID: message.Author.ID,
		UserID:      myID,
		Reason:      "test",
	}
	guild.Modlog.Log <- &modlog.CaseWarn{
		ModeratorID: message.Author.ID,
		UserID:      myID,
		Reason:      "test",
	}
	session.MessageReactionAdd(message.ChannelID, message.ID, "☑")
	return nil
}
