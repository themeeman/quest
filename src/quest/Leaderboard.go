package MyCommands

import (
	"github.com/bwmarrin/discordgo"
	commands ".././discordcommands"
	"sort"
	"fmt"
	"bytes"
	"sync"
)

type membersSorted []string
func (m membersSorted) Len() int           { return len(m) }
func (m membersSorted) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m membersSorted) Less(i, j int) bool { return members[m[i]].Experience > members[m[j]].Experience }

var members commands.Members

func Leaderboard(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string, bot commands.Bot) commands.BotError {
	guild := bot.Guilds.Get(commands.MustGetGuildID(session, message))
	members = guild.Members
	sorted := make(membersSorted, len(members))
	var i int
	for k := range members {
		sorted[i] = k
		i += 1
	}
	sort.Sort(sorted)
	var board []string
	if len(sorted) < 10 {
		board = sorted[:]
	} else {
		board = sorted[:10]
	}
	var buffer bytes.Buffer
	ms := getMembers(session, guild.ID, board)
	for i, id := range board {
		mem := guild.Members.Get(id)
		m := ms[i]
		if m != nil {
			if m.User.ID == message.Author.ID {
				buffer.WriteString(fmt.Sprintf("**%d. %s#%s - %d EXP**\n",
					i+1, m.User.Username, m.User.Discriminator, mem.Experience))
			} else {
				buffer.WriteString(fmt.Sprintf("%d. %s#%s - %d EXP\n",
					i+1, m.User.Username, m.User.Discriminator, mem.Experience))
			}
		} else {
			buffer.WriteString(fmt.Sprintf("%d. User Not Found - %d EXP\n",
				i+1, mem.Experience))
		}
	}
	index, _ := findIndex(sorted, message.Author.ID)
	if index > 9 {
		exp := guild.Members.Get(message.Author.ID).Experience
		buffer.WriteString(fmt.Sprintf("**%d. %s#%s - %d EXP\n**",
			index+1, message.Author.Username, message.Author.Discriminator, exp))
	}
	session.ChannelMessageSendEmbed(message.ChannelID, bot.Embed("Leaderboard", "", []*discordgo.MessageEmbedField{
		{
			Name:  "Ranks",
			Value: buffer.String(),
		},
	}))
	return nil
}

func findIndex(ss []string, s string) (int, bool) {
	for i, v := range ss {
		if v == s {
			return i, true
		}
	}
	return len(ss), false
}

func getMembers(session *discordgo.Session, guildID string, ids []string) []*discordgo.Member {
	var wg sync.WaitGroup
	wg.Add(len(ids))
	ms := make([]*discordgo.Member, len(ids))
	for i, id := range ids {
		go func(i int, id string) {
			ms[i], _ = session.GuildMember(guildID, id)
			wg.Done()
		}(i, id)
	}
	wg.Wait()
	return ms
}
