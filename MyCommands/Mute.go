package MyCommands

import (
	"github.com/bwmarrin/discordgo"
	commands "discordcommands"
	"strconv"
	"fmt"
	"strings"
	"time"
)

func Mute(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string, bot commands.Bot) commands.BotError {
	if len(args) >= 2 {
		ch, _ := session.State.Channel(message.ChannelID)
		dur, _ := strconv.Atoi(args["Minutes"])
		var user *discordgo.User
		if len(message.Mentions) == 0 {
			user = new(discordgo.User)
		} else {
			user = message.Mentions[0]
		}
		member, _ := session.State.Member(ch.GuildID, user.ID)
		g, _ := commands.FindGuildByID(bot.Guilds, ch.GuildID)
		for _, r := range member.Roles {
			if r == g.MuteRole.String {
				return commands.MutedError{
					Username:      user.Username,
					Discriminator: user.Discriminator,
				}
			}
		}
		err := session.GuildMemberRoleAdd(ch.GuildID, user.ID, g.MuteRole.String)
		if err != nil {
			fmt.Println(err)
			if strings.HasPrefix(err.Error(), "HTTP 403 Forbidden") {
				return commands.PermissionsError{}
			} else if strings.HasPrefix(err.Error(), "HTTP 400 Bad Request") {
				return commands.RoleError{ID: g.MuteRole.String}
			} else {
				return commands.UserNotFoundError{}
			}
		}
		m, _ := commands.FindMemberByID(g.Members, message.Author.ID)
		m.Mute.Time = time.Now()
		m.Mute.Valid = true
		m.MuteTime = dur
		go func() {
			time.Sleep(time.Second * time.Duration(dur))
			session.GuildMemberRoleRemove(ch.GuildID, user.ID, g.MuteRole.String)
		}()
		if args["Reason"] == "" {
			session.ChannelMessageSendEmbed(message.ChannelID, bot.Embed("Success!", fmt.Sprintf("Successfully muted %s#%s!", user.Username, user.Discriminator), nil))
		} else {
			session.ChannelMessageSendEmbed(message.ChannelID, bot.Embed("Success!", fmt.Sprintf("Successfully muted %s#%s! Reason: %s", user.Username, user.Discriminator, args["Reason"]), nil))
		}
	}
	return nil
}
