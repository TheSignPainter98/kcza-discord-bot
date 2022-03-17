package bot

import (
	"fmt"
	"golang-discord-bot/actions"
	"golang-discord-bot/cfg"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

const INITIAL_WATCHERS_SIZE = 10

type KczaBot struct {
	id             string
	session        *discordgo.Session
	watchers       []Watcher
	actionChannel  chan actions.DiscordAction
	botChannel     *string
	botChannelName string
}

func New(config *cfg.Config) (*KczaBot, error) {
	bot := &KczaBot{}
	bot.botChannelName = config.BotChannelName

	session, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return nil, err
	}
	bot.session = session

	u, err := bot.session.User("@me")
	if err != nil {
		return nil, err
	}

	bot.id = u.ID
	bot.session.AddHandler(bot.messageHandler)

	err = bot.session.Open()
	if err != nil {
		return nil, err
	}

	watchers := make([]Watcher, 0, INITIAL_WATCHERS_SIZE)
	web_watchers, err := InitWebWatchers(config, bot.GetBotChannel())
	if err != nil {
		return nil, err
	}
	bot.watchers = append(watchers, web_watchers...)
	bot.actionChannel = make(chan actions.DiscordAction)

	return bot, nil
}

func (bot KczaBot) HandleActions() {
	fmt.Println("Handling actions and there are ", len(bot.watchers), " of them")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	for _, w := range bot.watchers {
		go w.Watch(bot.actionChannel)
	}

	for {
		fmt.Println("Awaiting action...")
		select {
		case <-sig:
			return
		case action := <-bot.actionChannel:
			fmt.Println("Acting upon: " + fmt.Sprintf("%#v", action))
			action.Act(bot.session)
		}
	}
}

func (bot KczaBot) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == bot.id {
		return
	}

	if m.Content == "ping" {
		fmt.Println("Sending pong to " + m.ChannelID)
		s.ChannelMessageSend(m.ChannelID, "pong")
	}
}

func (bot KczaBot) GetBotChannel() string {
	if bot.botChannel == nil {
		for _, guild := range bot.session.State.Guilds {
			channels, _ := bot.session.GuildChannels(guild.ID)
			for _, c := range channels {
				if c.Type != discordgo.ChannelTypeGuildText {
					continue
				}
				if c.Name == bot.botChannelName {
					bot.botChannel = &c.ID
					goto end
				}
			}
		}
	}
end:
	return *bot.botChannel
}
