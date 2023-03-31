package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/itchyny/timefmt-go"
)

var timestampsFilePath = "timestamps.json"

type urlStorageStruct struct {
	Discord struct {
		WebhookUrl string
	}
	Bungie struct {
		Base string
		Rss  string
	}
	Twitter struct {
		Auth          string
		ApiBase       string
		TweetTemplate string
		QueryTemplate string
		Queries       struct {
			BungieHelp   string
			Destiny2Team string
		}
	}
}

type timestampStruct struct {
	Bungie              string
	TwitterBungieHelp   string
	TwitterDestiny2Team string
}

func getField(t *timestampStruct, field string) string {
	r := reflect.ValueOf(t)
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.String())
}

// Generate a timestamp of now in UTC
func timestamp() time.Time {
	ts := time.Now().UTC()
	DebugLogger.Println("Generated timestamp:", ts)
	return ts
}

// Convert an RFC1123 string into a time.Time object
func convertStrToTime(input string) time.Time {
	data, err := time.Parse(time.RFC1123, input)
	if err != nil {
		ErrorLogger.Println(err)
	}
	return data
}

// Convert a time.Time object into a RFC1123 string
func convertTimeToStr(input time.Time) string {
	return input.Format(time.RFC1123)
}

// Convert twitter's stupid timestamp to a time.Time object
func convertTwitterTimeStrToTime(input string) time.Time {
	// "Wed Mar 29 22:15:50 +0000 2023"
	// %d (0-padded) or %e (not padded) ?
	data, err := timefmt.Parse(input, "%a %b %d %H:%M:%S %z %Y")
	if err != nil {
		ErrorLogger.Println(err)
	}
	return data
}

// Reads in timestamps.json, creating a fresh one with the current timestamp if it doesn't exist
func readTimestampsFile() {
	data, err := os.ReadFile(timestampsFilePath)
	if err != nil {
		InfoLogger.Println("No existing timestamps file found; created one with current timestamp")
		newTime := timestamp()
		timestamps.Bungie = convertTimeToStr(newTime)
		timestamps.TwitterBungieHelp = convertTimeToStr(newTime)
		timestamps.TwitterDestiny2Team = convertTimeToStr(newTime)
		writeTimestampsFile()
	} else {
		InfoLogger.Println("Loaded timestamps.json")
		json.Unmarshal(data, &timestamps)
		DebugLogger.Println("Timestamp data read:", timestamps)
	}

}

// Writes to the timestamps.json file
func writeTimestampsFile() {
	data, err := json.Marshal(timestamps)
	if err != nil {
		ErrorLogger.Println(err)
	}

	err = os.WriteFile(timestampsFilePath, data, 0666)
	if err != nil {
		ErrorLogger.Println(err)
	}
	DebugLogger.Println("Timestamp data wrote:", timestamps)
}

// Fills out the various urls
func populateUrlStorage() {
	urls.Twitter.Auth = "https://api.twitter.com/1.1/guest/activate.json"
	urls.Twitter.ApiBase = "https://api.twitter.com/1.1/"
	urls.Twitter.TweetTemplate = "https://twitter.com/%s/status/%s"
	urls.Bungie.Base = "https://bungie.net"
	urls.Bungie.Rss = "https://www.bungie.net/en/rss/News"

	urls.Twitter.QueryTemplate = urls.Twitter.ApiBase + "statuses/user_timeline.json?screen_name=%s&exclude_replies=true&include_rts=false&count=50"
	urls.Twitter.Queries.BungieHelp = fmt.Sprintf(urls.Twitter.QueryTemplate, "BungieHelp")
	urls.Twitter.Queries.Destiny2Team = fmt.Sprintf(urls.Twitter.QueryTemplate, "Destiny2Team")

	DebugLogger.Println("UrlStorage:", urls)
}

func signalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Print("\r")
		DebugLogger.Println("Triggered signal handler")
		InfoLogger.Println("Goodbye")
		os.Exit(0)
	}()
}
