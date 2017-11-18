package sensors

import (
	"log"

	"github.com/mmcdole/gofeed"
)

type Headlines struct {
	Headlines []string `json:"headlines"`
}

var (
	fp *gofeed.Parser
)

func init() {
	fp = gofeed.NewParser()
}

func CheckHeadlines(feedURL string) []string {
	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		log.Fatal(err)
	}
	feedLines := []string{feed.Title}
	for _, item := range feed.Items {
		feedLines = append(feedLines, item.Title)
	}
	return feedLines
}
