package bot

import (
	"encoding/json"
	"fmt"
	"golang-discord-bot/actions"
	"golang-discord-bot/cfg"
	"io/ioutil"
)

type Watcher interface {
	Watch(chan actions.DiscordAction)
}

type WebWatcherSource struct {
	Name string `json: "name"`
	Url  string `json : "url"`
}

func InitWebWatchers(config *cfg.Config, outputChannel string) ([]Watcher, error) {
	fmt.Println("Reading from " + config.WebWatcherSource)
	rawWatchers, err := readWatchers(config.WebWatcherSource)
	if err != nil {
		return nil, err
	}

	watchers := make([]Watcher, len(rawWatchers))
	for i, r := range rawWatchers {
		watchers[i] = r.WebWatcher(outputChannel)
	}

	return watchers, nil
}

func readWatchers(loc string) ([]WebWatcherSource, error) {
	file, err := ioutil.ReadFile(loc)
	if err != nil {
		return nil, err
	}
	var watchers []WebWatcherSource
	if err = json.Unmarshal(file, &watchers); err != nil {
		return nil, err
	}
	return watchers, nil
}

func (w WebWatcherSource) WebWatcher(outputChannel string) Watcher {
	// TODO: parse urls, make this a factory
	return GeekHackWatcher{w.Name, w.Url, 123, outputChannel}
}
