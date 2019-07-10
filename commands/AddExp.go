package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/tomvanwoow/quest/modlog"
	"github.com/tomvanwoow/quest/utility"
	"github.com/tomvanwoow/quest/structures"
	"strconv"
	"strings"
)

func (bot *Bot) AddExp(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string) error {
	guild := bot.Guilds.Get(utility.MustGetGuildID(session, message))
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
	exp, _ := strconv.Atoi(strings.Replace(args["Value"], ",", "", -1))
	guild.Members.Apply(id, func(member *structures.Member) {
		member.Experience += int64(exp)
	})
	guild.RLock()
	defer guild.RUnlock()
	_ = session.MessageReactionAdd(message.ChannelID, message.ID, "☑")
	if guild.Modlog.Valid {
		guild.Modlog.Log <- &modlog.CaseAddExp{
			ModeratorID: message.Author.ID,
			Experience:  int64(exp),
			UserID:      id,
		}
	}
	return nil
}
