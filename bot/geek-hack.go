package bot

import (
	"errors"
	"fmt"
	"golang-discord-bot/actions"
	"golang-discord-bot/parse"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type TopicID int

type GeekHackWatcher struct {
	Name string
	Url  string
	TopicID
	outputChannel string
}

const retryBackoff = 1 * time.Minute

func (g GeekHackWatcher) Watch(act chan actions.DiscordAction) {
	fmt.Println("Hello, geekhack " + g.Name)
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

		act <- actions.MessageAction{g.outputChannel, fmt.Sprintf("Geekhack update for %v!", g.Name)}
		act <- actions.MessageAction{g.outputChannel, "This is another quite quick message"}
		act <- actions.LongMessage(g.outputChannel, bbc)
		time.Sleep(5 * time.Minute)
	}
}

func (g GeekHackWatcher) GetPageText() (string, error) {
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
