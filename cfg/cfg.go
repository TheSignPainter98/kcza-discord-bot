package cfg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type ConfigStruct struct {
	Token     string `json : "Token"`
	BotPrefix string `json : "BotPrefix"`
}

func ReadConfig() (*ConfigStruct, error) {
	fmt.Println("Reading config")
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		return nil, err
	}
	var cfg *ConfigStruct
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
