package vote

import "github.com/bwmarrin/discordgo"

func vote() {

}

func StartVote(s *discordgo.Session, channelID string, options []string) {
	// implementation of StartVote function
	if options[0] == "2" {
		embed := &discordgo.MessageEmbed{
			Author:      &discordgo.MessageEmbedAuthor{},
			Color:       0x00ff00,
			Description: "임베드 메세지",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Field1",
					Value:  "Value1",
					Inline: true,
				},
				{
					Name:   "Field2",
					Value:  "Value2",
					Inline: true,
				},
			},
			Title: "찬반 투표",
		}

		_, _ = s.ChannelMessageSendEmbed(channelID, embed)
	}
}
