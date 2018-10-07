package commands

import (
	commands "../../discordcommands"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"time"
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
		return UserNotFoundError{}
	}
	guild := bot.Guilds.Get(commands.MustGetGuildID(session, message))
	fmt.Println(guild.Cases, fmt.Sprintf("%p", &guild.Cases))
	g, err := session.State.Guild(commands.MustGetGuildID(session, message))
	if err != nil {
		return nil
	}
	var m *discordgo.Member
	if found, index := containsMember(g.Members, id); found {
		m = g.Members[index]
	} else {
		m, err = session.GuildMember(guild.ID, id)
		if err != nil {
			return UserNotFoundError{}
		}
	}
	if m.User.Bot {
		return fmt.Errorf("Cannot use `me` command on a bot!")
	}
	member := guild.Members.Get(id)
	rank := bot.UserGroup(session, g, m)
	s := []string{"Member", "Moderator", "Admin", "Owner"}
	title := fmt.Sprintf("User %s#%s", m.User.Username, m.User.Discriminator)
	fields := []*discordgo.MessageEmbedField{
		{
			Name:   "Experience",
			Value:  strconv.Itoa(int(member.Experience)),
			Inline: true,
		},
		{
			Name:   "Group",
			Value:  s[rank],
			Inline: true,
		},
	}
	now := time.Now().UTC()
	if member.MuteExpires.Valid && member.MuteExpires.Time.UTC().After(now) {
		dif := int(member.MuteExpires.Time.UTC().UnixNano() - now.UnixNano())
		hours := dif / int(time.Hour)
		minutes := (dif / int(time.Minute)) % 60
		seconds := (dif / int(time.Second)) % 60
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "Mute Time Left",
			Value: fmt.Sprintf("%d:%d:%d", hours, minutes, seconds),
		})
	}
	session.ChannelMessageSendEmbed(message.ChannelID, bot.Embed(title, "", fields))
	return nil
}
