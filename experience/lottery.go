package experience

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/tomvanwoow/quest/structures"
	"github.com/tomvanwoow/quest/utility"
	"math/rand"
)

func GrantLottery(session *discordgo.Session, guild *structures.Guild, member *structures.Member) {
	if guild.LotteryChance != 0 && rand.Intn(int(guild.LotteryChance)) == 0 {
		ch, err := session.UserChannelCreate(member.ID)
		u, _ := session.GuildMember(guild.ID, member.ID)
		var r int64
		if guild.LotteryLower > guild.LotteryUpper {
			r = int64(utility.RandInt(int(guild.LotteryUpper), int(guild.LotteryLower)))
		} else {
			r = int64(utility.RandInt(int(guild.LotteryLower), int(guild.LotteryUpper)))
		}
		if err == nil {
			guild, _ := session.Guild(guild.ID)
			if guild != nil {
				_, _ = session.ChannelMessageSend(ch.ID, fmt.Sprintf(`Looks like SOMEBODY is a lucky winner!
That's right, **%s#%s**, you won a grand total of %d Experience in **%s**! You should give yourself a pat on the back, you're a real winner in life!
🎉🎉🎉🎉🎉🎉🎉`, u.User.Username, u.User.Discriminator, r, guild.Name))
			} else {
				_, _ = session.ChannelMessageSend(ch.ID, fmt.Sprintf(`Looks like SOMEBODY is a lucky winner!
That's right, **%s#%s**, you won a grand total of %d Experience! You should give yourself a pat on the back, you're a real winner in life!
🎉🎉🎉🎉🎉🎉🎉`, u.User.Username, u.User.Discriminator, r))
			}
		}
		member.Experience += r
	}
}
