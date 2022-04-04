package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang-discord-bot/actions"
	"golang-discord-bot/cfg"
	"golang-discord-bot/parse"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type TryMakeWebWatcherFunc func(name, url, outputChannel string) Watcher

var webWatcherConstructors []TryMakeWebWatcherFunc

type Watcher interface {
	Watch(chan actions.DiscordAction)
}

type WebWatcherSource struct {
	Name string `json: "name"`
	Url  string `json : "url"`
}

type WebWatcher struct {
	Name          string
	Url           string
	OutputChannel string
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
	for _, f := range webWatcherConstructors {
		if o := f(w.Name, w.Url, outputChannel); o != nil {
			return o
		}
	}
	return nil
}

func (g WebWatcher) Watch(act chan actions.DiscordAction) {
	for true {
		fmt.Printf("Sleeping on %#v\n", g.Name)
		pageText, err := g.GetPageText()
		if err != nil {
			fmt.Println(err)
			time.Sleep(retryBackoff)
			continue
		}

		bbc, err := parse.BBCodeToMd(pageText)
		if err != nil {
			fmt.Println(err)
			time.Sleep(retryBackoff)
			continue
		}

		act <- actions.MessageAction{g.OutputChannel, fmt.Sprintf("Geekhack update for %v!", g.Name)}
		act <- actions.MessageAction{g.OutputChannel, "This is another quite quick message"}
		act <- actions.LongMessage(g.OutputChannel, bbc, 0)
		time.Sleep(5 * time.Minute)
	}
}

func (g WebWatcher) GetPageText() (string, error) {
	fmt.Printf("Getting page text from %s\n", g.Url)
	resp, err := http.Get(g.Url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("Returned status code %d\n", resp.StatusCode))
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	thing := doc.Find(".postarea").First()

	return strings.TrimSpace(thing.Text()), nil
}
