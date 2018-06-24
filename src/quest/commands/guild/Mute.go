package guild

import (
	"github.com/bwmarrin/discordgo"
	commands "../../../discordcommands"
	"strconv"
	"fmt"
	"strings"
	"time"
)

func Mute(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string, bot *commands.Bot) commands.BotError {
	ch, _ := session.State.Channel(message.ChannelID)
	var user *discordgo.User
	if len(args["User"]) == 18 {
		var err error
		user, err = session.User(args["User"])
		if err != nil {
			return commands.UserNotFoundError{}
		}
	} else if len(message.Mentions) > 0 {
		user = message.Mentions[0]
	} else {
		return commands.UserNotFoundError{}
	}
	member, _ := session.State.Member(ch.GuildID, user.ID)
	guild := bot.Guilds.Get(ch.GuildID)
	if !guild.MuteRole.Valid {
		return commands.MuteRoleError{}
	}
	for _, r := range member.Roles {
		if r == guild.MuteRole.String {
			return commands.MutedError{
				Username:      user.Username,
				Discriminator: user.Discriminator,
			}
		}
	}
	err := session.GuildMemberRoleAdd(ch.GuildID, user.ID, guild.MuteRole.String)
	if err != nil {
		fmt.Println(err)
		if strings.HasPrefix(err.Error(), "HTTP 403 Forbidden") {
			return commands.BotPermissionsError{}
		} else if strings.HasPrefix(err.Error(), "HTTP 400 Bad Request") {
			return commands.RoleError{ID: guild.MuteRole.String}
		} else {
			return commands.UserNotFoundError{}
		}
	}
	dur, _ := strconv.Atoi(strings.Replace(args["Minutes"], ",", "", -1))
	m := guild.Members.Get(message.Author.ID)
	m.MuteExpires.Time = time.Now().Add(time.Minute * time.Duration(dur))
	m.MuteExpires.Valid = true
	go func() {
		time.Sleep(time.Minute * time.Duration(dur))
		session.GuildMemberRoleRemove(ch.GuildID, user.ID, guild.MuteRole.String)
	}()
	if args["Reason"] == "" {
		session.ChannelMessageSendEmbed(message.ChannelID, bot.Embed("Success!", fmt.Sprintf("Successfully muted %s#%s!", user.Username, user.Discriminator), nil))
	} else {
		session.ChannelMessageSendEmbed(message.ChannelID, bot.Embed("Success!", fmt.Sprintf("Successfully muted %s#%s! Reason: %s", user.Username, user.Discriminator, args["Reason"]), nil))
	}
	return nil
}
