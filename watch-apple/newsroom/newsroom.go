package newsroom

import (
	"time"

	"github.com/mmcdole/gofeed"
)

// AppleNewsroomURL is url of Apple Newsroom RSS.
const AppleNewsroomURL = "https://www.apple.com/jp/newsroom/rss-feed.rss"

var now = time.Now()

// GetNewPosts checks if there are new posts on the Apple Newsroom and return them if exist.
func GetNewPosts(duration time.Duration) (posts []*gofeed.Item) {
	parser := gofeed.NewParser()
	feed, _ := parser.ParseURL(AppleNewsroomURL)
	for _, item := range feed.Items {
		updated, _ := time.Parse(time.RFC3339, item.Updated)
		if now.Sub(updated) < duration {
			posts = append(posts, item)
		}
	}
	return
}
