package guild

import (
	"github.com/bwmarrin/discordgo"
	commands "../../../discordcommands"
	"strings"
	"strconv"
	"log"
	"fmt"
	"math"
)

func Conv(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string, _ commands.Bot) commands.BotError {
	var number float64
	c := strings.Split(args["Number"], "e")
	v, err := strconv.ParseFloat(c[0], 32)
	if err != nil {
		log.Println(err)
		return nil
	}
	if len(c) == 1 {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%g", v))
		return nil
	}
	e, err := strconv.Atoi(c[1])
	if err != nil {
		log.Println(err)
		return nil
	}
	number = v * math.Pow10(e)
	session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%g", number))
	return nil
}
