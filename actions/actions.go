package actions

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

type DiscordAction interface {
	Act(*discordgo.Session) error
}

type ActionSeq []DiscordAction

type MessageAction struct {
	Channel string
	Message string
}

type TypingAction struct {
	Channel       string
	PauseDuration float64
}

func (seq ActionSeq) Act(s *discordgo.Session) error {
	for _, a := range seq {
		if e := a.Act(s); e != nil {
			return e
		}
	}
	return nil
}

func (msg MessageAction) Act(s *discordgo.Session) (err error) {
	if len(msg.Message) > 0 {
		fmt.Printf("Sending message of length %d to channel %s\n", len(msg.Message), msg.Channel)
		_, err = s.ChannelMessageSend(msg.Channel, msg.Message)
	}
	return
}

func LongMessage(channel string, msg string, typingDuration float64) ActionSeq {
	const blockSize = 2000
	const rollBackMax = blockSize / 10

	msgRunes := []rune(msg)
	n := len(msgRunes)
	maxMsgs := 1 + n/blockSize
	seq := make(ActionSeq, 0, 2*maxMsgs)
	typingDuration = typingDuration / float64(maxMsgs)

	for i := 0; i < n; i += blockSize {
		// Try to roll back
		intendedMax := intMin(i+blockSize, n-1)
		rollBack := 0
		for rollBack <= rollBackMax && !breakable(msgRunes[intendedMax-rollBack]) {
			rollBack++
		}
		intendedMax -= rollBack

		if typingDuration > 0 {
			seq = append(seq, TypingAction{channel, typingDuration})
		}
		seq = append(seq, MessageAction{channel, string(msgRunes[i:intendedMax])})
		i -= rollBack
	}

	return seq
}

func breakable(r rune) (canBreak bool) {
	breakables := map[rune]bool{
		' ':  true,
		'\t': true,
		'-':  true,
	}
	_, canBreak = breakables[r]
	return
}

func (t TypingAction) Act(s *discordgo.Session) error {
	if err := s.ChannelTyping(t.Channel); err != nil {
		return err
	}

	fmt.Printf("Using duration %f\n", t.PauseDuration)
	if t.PauseDuration > 0 {
		delay := time.Duration(int64(t.PauseDuration)) * time.Second
		time.Sleep(delay)
		fmt.Printf("Done waiting %#v\n", delay)
	}
	return nil
}

func intMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}
