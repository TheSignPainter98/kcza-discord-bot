package bot

import (
	"fmt"
	"regexp"
	"time"
)

type TopicID int

type GeekHackWatcher struct {
	WebWatcher
	TopicID
}

func init() {
	webWatcherConstructors = append(webWatcherConstructors, TryNewGeekHackWatcher)
}

const retryBackoff = 1 * time.Minute

func TryNewGeekHackWatcher(name, url, outputChannel string) Watcher {
	if !regexp.MustCompile(`geekhack\.org/index\.php`).MatchString(url) {
		return nil
	}

	rawTopic := regexp.MustCompile(`topic=\d+`).FindString(url)
	if len(rawTopic) == 0 {
		return nil
	}
	var topic TopicID
	fmt.Sscanf(url, "topic=%d", &topic)

	return &GeekHackWatcher{WebWatcher{name, url, outputChannel}, topic}
}
