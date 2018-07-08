package experience

import (
	"github.com/bwmarrin/discordgo"
	"time"
	"fmt"
	commands "../../discordcommands"
	"math/rand"
)

func GrantExp(bot *commands.Bot, session *discordgo.Session, message *discordgo.MessageCreate) {
	s := struct {
		Guild  string
		Member string
	}{
		Guild:  commands.MustGetGuildID(session, message),
		Member: message.Author.ID,
	}
	t, ok := bot.ExpTimes[s]
	g := bot.Guilds.Get(s.Guild)
	member := g.Members.Get(s.Member)
	if !ok || uint16(time.Since(t).Seconds()) > g.ExpReload {
		bot.ExpTimes[s] = time.Now()
		var r int64
		if g.ExpGainLower > g.ExpGainUpper {
			r = int64(rand.Intn(int(g.ExpGainLower+1-g.ExpGainUpper)) + int(g.ExpGainUpper))
		} else {
			r = int64(rand.Intn(int(g.ExpGainUpper+1-g.ExpGainLower)) + int(g.ExpGainLower))
		}
		member.Experience += r
		fmt.Println(s.Member, r)
		var a int
		if g.LotteryChance == 0 {
			a = 1
		} else {
			a = rand.Intn(int(g.LotteryChance))
		}
		fmt.Println(a)
		if a == 0 {
			ch, err := session.UserChannelCreate(s.Member)
			u, _ := session.GuildMember(s.Guild, s.Member)
			var r int64
			if g.LotteryLower > g.LotteryUpper {
				r = int64(rand.Intn(int(g.LotteryLower+1-g.LotteryUpper)) + int(g.LotteryUpper))
			} else {
				r = int64(rand.Intn(int(g.LotteryUpper+1-g.LotteryLower)) + int(g.LotteryLower))
			}
			if err == nil {
				guild, _ := session.Guild(s.Guild)
				if guild != nil {
					session.ChannelMessageSend(ch.ID, fmt.Sprintf(`Looks like SOMEBODY is a lucky winner!
That's right, **%s#%s**, you won a grand total of %d Experience in **%s**! You should give yourself a pat on the back, you're a real winner in life!
🎉🎉🎉🎉🎉🎉🎉`, u.User.Username, u.User.Discriminator, r, guild.Name))
				} else {
					session.ChannelMessageSend(ch.ID, fmt.Sprintf(`Looks like SOMEBODY is a lucky winner!
That's right, **%s#%s**, you won a grand total of %d Experience! You should give yourself a pat on the back, you're a real winner in life!
🎉🎉🎉🎉🎉🎉🎉`, u.User.Username, u.User.Discriminator, r))
				}
			}
			member.Experience += r
		}
	}
	GrantRoles(session, message, g, member)
}

