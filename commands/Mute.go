package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/tomvanwoow/quest/modlog"
)

func (bot *Bot) Mute(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string) error {
	ch, _ := session.State.Channel(message.ChannelID)
	var user *discordgo.User
	if len(args["User"]) == 18 {
		var err error
		user, err = session.User(args["User"])
		if err != nil {
			return UserNotFoundError{}
		}
	} else if len(message.Mentions) > 0 {
		user = message.Mentions[0]
	} else {
		return UserNotFoundError{}
	}
	member, _ := session.State.Member(ch.GuildID, user.ID)
	guild := bot.Guilds.Get(ch.GuildID)
	guild.RLock()
	defer guild.RUnlock()
	g, _ := session.Guild(ch.GuildID)
	if bot.UserGroup(session, g, member) >= PermissionAdmin {
		return fmt.Errorf("That user is a admin, I can't mute them!")
	} else if member.User.Bot {
		return fmt.Errorf("Can't mute a bot")
	}
	if !guild.MuteRole.Valid {
		return fmt.Errorf("No mute role has been configured for the server! Use q:set muterole [Value]")
	}
	for _, r := range member.Roles {
		if r == guild.MuteRole.String {
			return MutedError{
				Username:      user.Username,
				Discriminator: user.Discriminator,
			}
		}
	}
	err := session.GuildMemberRoleAdd(ch.GuildID, user.ID, guild.MuteRole.String)
	if err != nil {
		fmt.Println(err)
		if strings.HasPrefix(err.Error(), "HTTP 403 Forbidden") {
			return fmt.Errorf("Make sure the bot has Manage Roles Permission in Discord!")
		} else if strings.HasPrefix(err.Error(), "HTTP 400 Bad Request") {
			return RoleError{ID: guild.MuteRole.String}
		} else {
			return UserNotFoundError{}
		}
	}
	dur, _ := strconv.Atoi(strings.Replace(args["Minutes"], ",", "", -1))
	m := guild.Members.Get(user.ID)
	m.MuteExpires.Time = time.Now().UTC().Add(time.Minute * time.Duration(dur))
	m.MuteExpires.Valid = true
	go func() {
		time.Sleep(time.Minute * time.Duration(dur))
		session.GuildMemberRoleRemove(ch.GuildID, user.ID, guild.MuteRole.String)
		err := session.GuildMemberRoleRemove(guild.ID, user.ID, guild.MuteRole.String)
		if err == nil && guild.Modlog.Valid {
			guild.Modlog.Log <- &modlog.CaseUnmute{
				ModeratorID: "412702549645328397",
				UserID:      user.ID,
				Reason:      "Auto",
			}
		}
	}()
	if args["Reason"] == "" {
		session.ChannelMessageSendEmbed(message.ChannelID, bot.Embed("Success!", fmt.Sprintf("Successfully muted %s#%s!", user.Username, user.Discriminator), nil))
	} else {
		session.ChannelMessageSendEmbed(message.ChannelID, bot.Embed("Success!", fmt.Sprintf("Successfully muted %s#%s! Reason: %s", user.Username, user.Discriminator, args["Reason"]), nil))
	}
	if guild.Modlog.Valid {
		guild.Modlog.Log <- &modlog.CaseMute{
			ModeratorID: message.Author.ID,
			UserID:      user.ID,
			Duration:    dur,
			Reason:      args["Reason"],
		}
	}
	return nil
}
