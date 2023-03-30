package main

import (
	"github.com/mmcdole/gofeed"
)

// Attempts to return a *gofeed.Feed (up to 3 tries)
func getBungieFeed(parser *gofeed.Parser) (*gofeed.Feed, error) {
	tries := 1
	var feed *gofeed.Feed
	var err error = nil

	for tries <= 3 {
		feed, err = parser.ParseURL(urls.Bungie.Rss)
		if err != nil {
			tries++
		} else {
			break
		}
	}
	return feed, err

}

// Checks for new bungie.net posts, sending a webhook message if there's a new one
func parseBungiePosts(parser *gofeed.Parser) bool {
	feed, err := getBungieFeed(parser)
	if err != nil {
		ErrorLogger.Println("Failed to get feed")
		return false
	}

	newestItem := feed.Items[0]

	lastAlert := convertStrToTime(timestamps.Bungie)
	latestPost := convertStrToTime(newestItem.Published)

	if latestPost.Before(lastAlert) || latestPost.Equal(lastAlert) {
		InfoLogger.Println("Up to date: Bungie.net")
		return false
	} else {
		AlertLogger.Println("New Bungie.net post")
		content := newestItem.Title + "\n"
		content += urls.Bungie.Base + newestItem.Link
		sendDiscordWebhook(content)
		timestamps.Bungie = convertTimeToStr(latestPost)
		return true
	}

}
