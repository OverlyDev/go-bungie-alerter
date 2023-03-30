package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/mmcdole/gofeed"
)

var (
	InfoLogger   *log.Logger
	ErrorLogger  *log.Logger
	AlertLogger  *log.Logger
	urls         url_storage_struct
	twitter_auth twitter_auth_struct
	timestamps   timestamp_struct
)

func init() {
	// Set up the loggers
	InfoLogger = log.New(log.Default().Writer(), "INFO | ", log.Ldate|log.Ltime|log.Lshortfile|log.LUTC)
	ErrorLogger = log.New(log.Default().Writer(), "ERROR | ", log.Ldate|log.Ltime|log.Lshortfile|log.LUTC)
	AlertLogger = log.New(log.Default().Writer(), "ALERT | ", log.Ldate|log.Ltime|log.Lshortfile|log.LUTC)

	// Get the discord webhook url from env, failing if not found
	godotenv.Load()
	_, found := os.LookupEnv("DISCORD_WEBHOOK")
	if !found {
		ErrorLogger.Fatalln("DISCORD_WEBHOOK not found in evnironment, exiting")
	} else {
		urls.Discord.WebhookUrl = os.Getenv("DISCORD_WEBHOOK")
	}

	// Set up all the urls
	populate_url_storage()

	// Twitter auth setup stuff
	get_twitter_auth()
	if twitter_auth.Guest == "" {
		ErrorLogger.Fatalln("Failed to obtain twitter auth token, exiting")
	}

	// Load timestamps file
	read_timestamps_file()
}

func main() {
	feed_parser := gofeed.NewParser()

	InfoLogger.Println("Starting")
	for {
		InfoLogger.Println("Getting Bungie.net feed")
		changes := false || parse_bungie_posts(feed_parser)

		InfoLogger.Println("Getting tweets")
		changes = changes || check_for_tweets()

		if changes {
			AlertLogger.Println("Changes to timestamps, writing to disk")
			write_timestamps_file()
		} else {
			InfoLogger.Println("No changes to timestamps")
		}
		InfoLogger.Println("Sleeping 60s")
		time.Sleep(60 * time.Second)
	}

}
