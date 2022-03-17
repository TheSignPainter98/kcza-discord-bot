package main

import (
	"fmt"
	"golang-discord-bot/bot"
	"golang-discord-bot/cfg"
)

func main() {
	conf, err := cfg.ReadConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	kczaBot, err := bot.New(conf)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	kczaBot.HandleActions()
}
