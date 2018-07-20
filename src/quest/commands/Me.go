package commands

import (
	"github.com/bwmarrin/discordgo"
	commands "../../discordcommands"
	"strconv"
	"fmt"
	"time"
	"../permissions"
)

func (bot *Bot) Me(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string) error {
	var id string
	if args["User"] == "" {
		id = message.Author.ID
	} else if len(args["User"]) == 18 {
		id = args["User"]
	} else if len(message.Mentions) > 0 {
		id = message.Mentions[0].ID
	} else {
		return commands.UserNotFoundError{}
	}
	guild := bot.Guilds.Get(commands.MustGetGuildID(session, message))
	m, err := session.GuildMember(guild.ID, id)
	if err != nil {
		return commands.UserNotFoundError{}
	}
	if m.User.Bot {
		return commands.CustomError("Cannot use `me` command on a bot!")
	}
	member := guild.Members.Get(id)
	g, _ := session.Guild(commands.MustGetGuildID(session, message))
	rank := permissions.GetPermissionLevel(session, m, guild, g.OwnerID)
	s := []string{"Member", "Moderator", "Admin", "Owner"}
	title := fmt.Sprintf("User %s#%s", m.User.Username, m.User.Discriminator)
	fields := []*discordgo.MessageEmbedField{
		{
			Name:  "Experience",
			Value: strconv.Itoa(int(member.Experience)),
		},
		{
			Name:  "Group",
			Value: s[rank],
		},
	}
	now := time.Now().UTC()
	if member.MuteExpires.Valid && member.MuteExpires.Time.UTC().After(now) {
		dif := member.MuteExpires.Time.UTC().UnixNano() - now.UnixNano()
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "Mute Time Left (Seconds)",
			Value: strconv.Itoa(int(dif / int64(time.Second))),
		})
	}
	session.ChannelMessageSendEmbed(message.ChannelID, bot.Embed(title, "", fields))
	return nil
}
