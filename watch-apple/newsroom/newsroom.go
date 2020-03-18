package newsroom

import (
	"time"

	"github.com/mmcdole/gofeed"
)

// AppleNewsroomURL is url of Apple Newsroom RSS.
const AppleNewsroomURL = "https://www.apple.com/jp/newsroom/rss-feed.rss"

var now = time.Now()

// GetPosts gets news from Apple Newsroom
func GetPosts(after time.Time) (posts []*gofeed.Item) {
	parser := gofeed.NewParser()
	feed, _ := parser.ParseURL(AppleNewsroomURL)
	for _, item := range feed.Items {
		updated, _ := time.Parse(time.RFC3339, item.Updated)
		if updated.After(after) {
			posts = append(posts, item)
		}
	}
	return
}
