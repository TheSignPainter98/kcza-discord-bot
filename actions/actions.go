package actions

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type DiscordAction interface {
	Act(*discordgo.Session)
}

type ActionSeq []DiscordAction

type MessageAction struct {
	Channel string
	Message string
}

func (seq ActionSeq) Act(s *discordgo.Session) {
	for _, a := range seq {
		a.Act(s)
	}
}

func (msg MessageAction) Act(s *discordgo.Session) {
	fmt.Printf("Sending message of length %d to channel %s: %#v\n", len(msg.Message), msg.Channel, msg.Message)
	s.ChannelMessageSend(msg.Channel, msg.Message)
}

func LongMessage(channel string, msg string) ActionSeq {
	n := len(msg)
	const blockSize = 2000

	seq := make(ActionSeq, 0, 1+n/blockSize)

	for i := 0; i < n; i += blockSize {
		seq = append(seq, MessageAction{channel, string(msg[i:intMin(i+blockSize, n-1)])})
	}

	return seq
}

func intMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}