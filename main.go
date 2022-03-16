package main

import (
	"fmt"
	"golang-discord-bot/bot"
	"golang-discord-bot/cfg"
	"os"
	"os/signal"
)

func main() {
	conf, err := cfg.ReadConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bot.Start(conf)

	fmt.Println("Bot running, hanging...")
	hang := make(chan os.Signal, 1)
	signal.Notify(hang, os.Interrupt)
	<-hang
}
