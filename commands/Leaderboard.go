package commands

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/tomvanwoow/quest/structures"
	"github.com/tomvanwoow/quest/utility"
)

func (bot *Bot) Leaderboard(session *discordgo.Session, message *discordgo.MessageCreate, _ map[string]string) error {
	guildID := utility.MustGetGuildID(session, message)
	leaderBoard, err := structures.GetTopMembers(bot.DB, guildID, 10)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, "```Internal server error - try again later```")
		return nil
	}
	ids := make([]string, len(leaderBoard))
	for i, m := range leaderBoard {
		ids[i] = m.ID
	}
	members := getMembers(session, guildID, ids)
	var buffer bytes.Buffer
	guild := bot.Guilds.Get(guildID)
	for index, m := range members {
		guild.RLock()
		mem := guild.Members.Get(m.ID)
		guild.RUnlock()
		if m.Member != nil {
			if m.User.ID == message.Author.ID {
				buffer.WriteString(fmt.Sprintf("**%d. %s#%s - %d EXP**\n",
					index+1, m.User.Username, m.User.Discriminator, mem.Experience))
			} else {
				buffer.WriteString(fmt.Sprintf("%d. %s#%s - %d EXP\n",
					index+1, m.User.Username, m.User.Discriminator, mem.Experience))
			}
		} else {
			buffer.WriteString(fmt.Sprintf("%d. *User Not Found* - %d EXP\n",
				index+1, mem.Experience))
		}
	}
	specialIndex, _ := findIndex(sorted.IDs, message.Author.ID)
	if specialIndex > 9 {
		guild.RLock()
		mem := guild.Members.Get(sorted.IDs[specialIndex])
		guild.RUnlock()
		buffer.WriteString(fmt.Sprintf("**%d. %s#%s - %d EXP\n**",
			specialIndex+1, message.Author.Username, message.Author.Discriminator, mem.Experience))
	}
	session.ChannelMessageSendEmbed(message.ChannelID, bot.Embed("Leaderboard", buffer.String(), nil))
	return nil
}

type memberData struct {
	Index int
	ID    string
	*discordgo.Member
}

func getMembers(session *discordgo.Session, guildID string, ids []string) []memberData {
	var wg sync.WaitGroup
	wg.Add(len(ids))
	ms := make([]memberData, len(ids))
	guild, _ := session.State.Guild(guildID)
	if guild == nil {
		return nil
	}
	for i, id := range ids {
		go func(i int, id string) {
			if found, index := containsMember(guild.Members, id); found {
				ms[i] = memberData{
					Index:  i,
					ID:     id,
					Member: guild.Members[index],
				}
			} else {
				m, _ := session.State.Member(guildID, id)
				if m != nil {
					session.State.MemberAdd(m)
				}
				ms[i] = memberData{
					Index:  i,
					ID:     id,
					Member: m,
				}
			}
			wg.Done()
		}(i, id)
	}
	wg.Wait()
	return ms
}

func containsMember(members []*discordgo.Member, id string) (found bool, index int) {
	for i, v := range members {
		if v.User.ID == id {
			return true, i
		}
	}
	return false, 0
}
