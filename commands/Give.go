package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/tomvanwoow/quest/utility"
	"strconv"
)

func (bot *Bot) Give(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string) error {
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
	item, _ := strconv.Atoi(args["Item"])
	amount, _ := strconv.Atoi(args["Amount"])
	member := bot.Guilds.Get(utility.MustGetGuildID(session, message)).Members.Get(id)
	member.Lock()
	member.Chests[uint(item)] += uint(amount)
	member.Unlock()
	session.MessageReactionAdd(message.ChannelID, message.ID, "☑")
	return nil
}
