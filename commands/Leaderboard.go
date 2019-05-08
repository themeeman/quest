package commands

import (
	"github.com/tomvanwoow/quest/structures"
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/tomvanwoow/quest/utility"
	"sort"
	"sync"
)

type membersSorted struct {
	IDs []string
	structures.Members
}

func (m membersSorted) Len() int      { return len(m.IDs) }
func (m membersSorted) Swap(i, j int) { m.IDs[i], m.IDs[j] = m.IDs[j], m.IDs[i] }
func (m membersSorted) Less(i, j int) bool {
	return m.Members[m.IDs[i]].Experience > m.Members[m.IDs[j]].Experience
}

func (bot *Bot) Leaderboard(session *discordgo.Session, message *discordgo.MessageCreate, _ map[string]string) error {
	guild := bot.Guilds.Get(utility.MustGetGuildID(session, message))
	sorted := membersSorted{
		IDs:     make([]string, len(guild.Members)),
		Members: guild.Members,
	}

	{
		var i int
		for id := range guild.Members {
			sorted.IDs[i] = id
			i += 1
		}
	}

	sort.Sort(sorted)
	var board []string
	if len(sorted.IDs) < 10 {
		board = sorted.IDs[:]
	} else {
		board = sorted.IDs[:10]
	}
	members := getMembers(session, guild.ID, board)
	var buffer bytes.Buffer
	for index, m := range members {
		mem := guild.Members.Get(m.ID)
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
		mem := guild.Members.Get(sorted.IDs[specialIndex])
		buffer.WriteString(fmt.Sprintf("**%d. %s#%s - %d EXP\n**",
			specialIndex+1, message.Author.Username, message.Author.Discriminator, mem.Experience))
	}
	session.ChannelMessageSendEmbed(message.ChannelID, bot.Embed("Leaderboard", buffer.String(), nil))
	return nil
}

func findIndex(ss []string, s string) (int, bool) {
	for i, v := range ss {
		if v == s {
			return i, true
		}
	}
	return len(ss) - 1, false
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
