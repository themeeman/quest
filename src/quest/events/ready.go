package events

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"fmt"
	commands "../../discordcommands"
	"time"
)

func Ready(bot *commands.Bot) func(*discordgo.Session, *discordgo.Ready) {
	return func(session *discordgo.Session, ready *discordgo.Ready) {
		var err error
		guilds, err := commands.QueryAllData(bot.DB)
		if err != nil {
			log.Println("b", err)
		}
		for _, v := range guilds {
			fmt.Println(v)
		}
		bot.Guilds = guilds
		session.UpdateStatus(0, "q:help")
		go func() {
			for {
				time.Sleep(time.Minute * 10)
				err := commands.PostAllData(bot.DB, bot.Guilds)
				if err != nil {
					log.Println(err)
				} else {
					log.Println("Successfully commited all data")
				}
			}
		}()
	}
}