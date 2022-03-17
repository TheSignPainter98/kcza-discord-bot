package cfg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Token            string `json : "Token"`
	BotPrefix        string `json : "BotPrefix"`
	WebWatcherSource string `json : "WebWatchers"`
	BotChannelName   string `json : "BotChannelName"`
}

func ReadConfig() (*Config, error) {
	fmt.Println("Reading config")
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		return nil, err
	}
	var cfg *Config
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
