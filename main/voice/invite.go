package voice

import "github.com/bwmarrin/discordgo"

func Invite(s *discordgo.Session, m *discordgo.MessageCreate) {
	guild, _ := s.State.Guild(m.GuildID)
	var vs *discordgo.VoiceState

	for _, voiceState := range guild.VoiceStates {
		if voiceState.UserID == m.Author.ID {
			vs = voiceState
			break
		}
	}

	if vs != nil {
		_, _ = s.ChannelVoiceJoin(m.GuildID, vs.ChannelID, false, false)
	} else {
		_, _ = s.ChannelMessageSend(m.ChannelID, "음악을 재생하려면 음성채널에 속해 있어야 합니다.")
	}
}
