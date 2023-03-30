package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/mmcdole/gofeed"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	AlertLogger *log.Logger
	urls        urlStorageStruct
	twitterAuth twitterAuthStruct
	timestamps  timestampStruct
)

func init() {
	printVersion()

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
	populateUrlStorage()

	// Twitter auth setup stuff
	getTwitterAuth()
	if twitterAuth.Guest == "" {
		ErrorLogger.Fatalln("Failed to obtain twitter auth token, exiting")
	}

	// Load timestamps file
	readTimestampsFile()
}

func main() {
	feedParser := gofeed.NewParser()

	InfoLogger.Println("Starting")
	for {
		changes := false

		InfoLogger.Println("Getting Bungie.net feed")
		changes = changes || parseBungiePosts(feedParser)

		InfoLogger.Println("Getting tweets")
		changes = changes || checkForTweets()

		if changes {
			AlertLogger.Println("Changes to timestamps, writing to disk")
			writeTimestampsFile()
		} else {
			InfoLogger.Println("No changes to timestamps")
		}
		InfoLogger.Println("Sleeping 60s")
		time.Sleep(60 * time.Second)
	}

}
