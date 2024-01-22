package alarm

import (
	"fmt"
	customTime "gelato/main/time"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"time"
)

var (
	Quit chan struct{}
)

var (
	IsHeraldOn bool
)

func SendAlarm(s *discordgo.Session, channelID, minutes string) {
	n, _ := strconv.Atoi(minutes)

	_, _ = s.ChannelMessageSend(channelID, fmt.Sprintf("%s\n%d분 후에 메시지를 보냅니다.", customTime.PrintCurrentTime(), n))
	<-time.After(time.Duration(n) * time.Minute)
	_, _ = s.ChannelMessageSend(channelID, fmt.Sprintf("%s\n%d분이 지났습니다.", customTime.PrintCurrentTime(), n))
}

func SendDailyMessage(s *discordgo.Session, channelID string, quit <-chan struct{}) {
	location, _ := time.LoadLocation("Asia/Seoul")

	for {
		now := time.Now().In(location)
		next := now

		if now.Hour() < 2 {
			next = time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, now.Location())
		} else {
			next = now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), 2, 0, 0, 0, next.Location())
		}

		duration := next.Sub(now)
		fmt.Println(duration)

		time.Sleep(duration)

		select {
		case <-time.After(duration):
			_, _ = s.ChannelMessageSend(channelID, "새벽 2시입니다. 모두 나가!\n안나가? 얼른나가!")
		case <-quit:
			return
		}
	}
}
