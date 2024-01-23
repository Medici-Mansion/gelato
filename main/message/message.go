package message

import (
	"gelato/main/alarm"
	customtime "gelato/main/time"
	"gelato/main/voice"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func Create(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	cmd := strings.Split(m.Content, " ")

	if m.Content == "!i" {
		helpMessage := "### Gelato 사용 방법\n" +
			"```\n" +
			"!i\n" +
			"- 젤라또의 사용 방법에 대해서 알려줍니다.\n\n" +
			"!h on/off\n" +
			"- 젤라또가 오전 2시 감시를 시작/종료합니다.\n\n" +
			"!a 숫자\n" +
			"- 입력한 숫자만큼의 시간(분) 뒤에 메시지를 보내줍니다.\n\n" +
			"!t\n" +
			"- 현재 시각을 알려줍니다.\n\n" +
			"```"

		_, _ = s.ChannelMessageSend(m.ChannelID, helpMessage)
	}

	if m.Content == "!t" {
		_, _ = s.ChannelMessageSend(m.ChannelID, customtime.PrintCurrentTime())
		return
	}

	if len(cmd) > 1 && cmd[0] == "!a" {
		go alarm.SendAlarm(s, m.ChannelID, cmd[1])
	}

	if len(cmd) > 1 && cmd[0] == "!h" && cmd[1] == "on" {
		if alarm.IsHeraldOn {
			_, _ = s.ChannelMessageSend(m.ChannelID, "Gelato가 이미 감시중입니다.")
			return
		}

		_, _ = s.ChannelMessageSend(m.ChannelID, "Gelato가 당신들을 감시합니다.")
		alarm.IsHeraldOn = true
		alarm.Quit = make(chan struct{})
		go alarm.SendDailyMessage(s, m.ChannelID, alarm.Quit)
	}

	if len(cmd) > 1 && cmd[0] == "!h" && cmd[1] == "off" {
		if !alarm.IsHeraldOn {
			_, _ = s.ChannelMessageSend(m.ChannelID, "Gelato가 감시중이지 않습니다.")
		}

		if alarm.Quit != nil {
			_, _ = s.ChannelMessageSend(m.ChannelID, "Gelato가 감시를 종료합니다.")
			close(alarm.Quit)
			alarm.Quit = nil
			alarm.IsHeraldOn = false
		}
	}

	if len(cmd) > 1 && cmd[0] == "!m" {
		voice.Invite(s, m)
	}
}
