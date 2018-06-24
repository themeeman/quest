package guild

import (
	"github.com/bwmarrin/discordgo"
	commands "../../../discordcommands"
	"sort"
	"fmt"
	"bytes"
	"sync"
)

type membersSorted struct {
	IDs []string
	commands.Members
}

func (m membersSorted) Len() int      { return len(m.IDs) }
func (m membersSorted) Swap(i, j int) { m.IDs[i], m.IDs[j] = m.IDs[j], m.IDs[i] }
func (m membersSorted) Less(i, j int) bool {
	return m.Members[m.IDs[i]].Experience > m.Members[m.IDs[j]].Experience
}

func Leaderboard(session *discordgo.Session, message *discordgo.MessageCreate, _ map[string]string, bot *commands.Bot) commands.BotError {
	guild := bot.Guilds.Get(commands.MustGetGuildID(session, message))
	sorted := membersSorted{
		IDs:     make([]string, len(guild.Members)),
		Members: guild.Members,
	}
	var i int
	for k := range guild.Members {
		sorted.IDs[i] = k
		i += 1
	}
	sort.Sort(sorted)
	var board []string
	if len(sorted.IDs) < 10 {
		board = sorted.IDs[:]
	} else {
		board = sorted.IDs[:10]
	}
	embed, _ := session.ChannelMessageSendEmbed(message.ChannelID, bot.Embed("Leaderboard", "", nil))
	if embed == nil {
		return nil
	}
	ch := make(chan *memberData)
	ms := make([]*memberData, len(board))
	go getMembers(session, guild.ID, board, ch)
	for {
		go func(ms []*memberData) {
			var buffer bytes.Buffer
			for index, m := range ms {
				if m != nil {
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
						buffer.WriteString(fmt.Sprintf("%d. User Not Found - %d EXP\n",
							index+1, mem.Experience))
					}
				} else {
					buffer.WriteString(fmt.Sprintf("%d. Loading...\n", index+1))
				}
			}
			specialIndex, _ := findIndex(sorted.IDs, message.Author.ID)
			if specialIndex > 9 {
				mem := guild.Members.Get(sorted.IDs[specialIndex])
				buffer.WriteString(fmt.Sprintf("**%d. %s#%s - %d EXP\n**",
					specialIndex+1, message.Author.Username, message.Author.Discriminator, mem.Experience))
			}
			session.ChannelMessageEditEmbed(message.ChannelID, embed.ID, bot.Embed("Leaderboard", "", []*discordgo.MessageEmbedField{
				{
					Name:  "Ranks",
					Value: buffer.String(),
				},
			}))
		}(ms)
		member := <-ch
		if member == nil {
			break
		}
		ms[member.Index] = member
	}
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

type memberData struct {
	Index int
	ID    string
	*discordgo.Member
}

func getMembers(session *discordgo.Session, guildID string, ids []string, ch chan *memberData) []*discordgo.Member {
	var wg sync.WaitGroup
	wg.Add(len(ids))
	ms := make([]*discordgo.Member, len(ids))
	for i, id := range ids {
		go func(i int, id string) {
			ms[i], _ = session.GuildMember(guildID, id)
			if ch != nil {
				ch <- &memberData{
					Index:  i,
					ID:     id,
					Member: ms[i],
				}
			}
			wg.Done()
		}(i, id)
	}
	wg.Wait()
	ch <- nil
	return ms
}
