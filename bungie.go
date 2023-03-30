package main

import (
	"github.com/mmcdole/gofeed"
)

// Attempts to return a *gofeed.Feed (up to 3 tries)
func get_feed(parser *gofeed.Parser) (*gofeed.Feed, error) {
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
func parse_bungie_posts(parser *gofeed.Parser) bool {
	feed, err := get_feed(parser)
	if err != nil {
		ErrorLogger.Println("Failed to get feed")
		return false
	}

	newest_item := feed.Items[0]

	last_alert := convert_RFC1123_str_to_time(timestamps.Bungie)
	latest_post := convert_RFC1123_str_to_time(newest_item.Published)

	if latest_post.Before(last_alert) || latest_post.Equal(last_alert) {
		InfoLogger.Println("Up to date: Bungie.net")
		return false
	} else {
		AlertLogger.Println("New Bungie.net post")
		content := newest_item.Title + "\n"
		content += urls.Bungie.Base + newest_item.Link
		send_discord_webhook(content)
		timestamps.Bungie = convert_time_to_RFC1123_str(latest_post)
		return true
	}

}
