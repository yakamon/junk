package newsroom

import (
	"time"

	"github.com/mmcdole/gofeed"
)

const (
	// AppleNewsroomRssURL is URL of apple.com newsroom rss.
	AppleNewsroomRssURL = "https://www.apple.com/jp/newsroom/rss-feed.rss"
)

// GetNewPosts checks the Apple Newsroom for new posts and return them if exist.
func GetNewPosts(duration time.Duration) []*gofeed.Item {
	now := time.Now()
	parser := gofeed.NewParser()
	feed, _ := parser.ParseURL(AppleNewsroomRssURL)

	var newPosts []*gofeed.Item
	for _, item := range feed.Items {
		updated, _ := time.Parse(time.RFC3339, item.Updated)
		if now.Sub(updated) < duration {
			newPosts = append(newPosts, item)
		}
	}

	return newPosts
}
