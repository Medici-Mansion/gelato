package message

import (
	"gelato/main/alarm"
	customtime "gelato/main/time"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func Create(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	cmd := strings.Split(m.Content, " ")

	if m.Content == "time" {
		_, _ = s.ChannelMessageSend(m.ChannelID, customtime.PrintCurrentTime())
		return
	}

	if len(cmd) > 1 && cmd[0] == "alarm" {
		go alarm.SendAlarm(s, m.ChannelID, cmd[1])
	}

	if len(cmd) > 1 && cmd[0] == "gelato" && cmd[1] == "on" {
		if alarm.IsHeraldOn {
			_, _ = s.ChannelMessageSend(m.ChannelID, "gelato가 이미 실행중입니다.")
			return
		}

		_, _ = s.ChannelMessageSend(m.ChannelID, "gelato를 실행합니다.")
		alarm.IsHeraldOn = true
		alarm.Quit = make(chan struct{})
		go alarm.SendDailyMessage(s, m.ChannelID, alarm.Quit)
	}

	if len(cmd) > 1 && cmd[0] == "gelato" && cmd[1] == "off" {
		if !alarm.IsHeraldOn {
			_, _ = s.ChannelMessageSend(m.ChannelID, "gelato가 실행중이지 않습니다.")
		}

		if alarm.Quit != nil {
			_, _ = s.ChannelMessageSend(m.ChannelID, "gelato를 종료합니다.")
			close(alarm.Quit)
			alarm.Quit = nil
			alarm.IsHeraldOn = false
		}
	}
}
