package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/itchyny/timefmt-go"
)

var timestamps_file_path = "timestamps.json"

type url_storage_struct struct {
	Discord struct {
		WebhookUrl string
	}
	Bungie struct {
		Base string
		Rss  string
	}
	Twitter struct {
		Auth    string
		ApiBase string
		Queries struct {
			BungieHelp   string
			Destiny2Team string
		}
	}
}

type timestamp_struct struct {
	Bungie              string
	TwitterBungieHelp   string
	TwitterDestiny2Team string
}

// Generate a timestamp of now in UTC
func timestamp() time.Time {
	return time.Now().UTC()
}

// Convert an RFC1123 string into a time.Time object
func convert_RFC1123_str_to_time(input string) time.Time {
	data, err := time.Parse(time.RFC1123, input)
	if err != nil {
		fmt.Println(err)
	}
	return data
}

// Convert a time.Time object into a RFC1123 string
func convert_time_to_RFC1123_str(input time.Time) string {
	return input.Format(time.RFC1123)
}

func convert_twitter_time_str_to_time(input string) time.Time {
	// "Wed Mar 29 22:15:50 +0000 2023"
	// %d (0-padded) or %e (not padded) ?
	data, err := timefmt.Parse(input, "%a %b %d %H:%M:%S %z %Y")
	if err != nil {
		fmt.Println(err)
	}
	return data
}

func read_timestamps_file() {
	data, err := os.ReadFile(timestamps_file_path)
	if err != nil {
		InfoLogger.Println("No existing timestamps file found; created one with current timestamp")
		new_time := timestamp()
		timestamps.Bungie = convert_time_to_RFC1123_str(new_time)
		timestamps.TwitterBungieHelp = convert_time_to_RFC1123_str(new_time)
		timestamps.TwitterDestiny2Team = convert_time_to_RFC1123_str(new_time)
		write_timestamps_file()
	} else {
		InfoLogger.Println("Loaded timestamps.json")
		json.Unmarshal(data, &timestamps)
	}
}

func write_timestamps_file() {
	data, err := json.Marshal(timestamps)
	if err != nil {
		ErrorLogger.Println(err)
	}

	err = os.WriteFile(timestamps_file_path, data, 0666)
	if err != nil {
		ErrorLogger.Println(err)
	}
}

func populate_url_storage() {
	urls.Twitter.Auth = "https://api.twitter.com/1.1/guest/activate.json"
	urls.Twitter.ApiBase = "https://api.twitter.com/1.1/"
	urls.Bungie.Base = "https://bungie.net"
	urls.Bungie.Rss = "https://www.bungie.net/en/rss/News"

	query_template := urls.Twitter.ApiBase + "statuses/user_timeline.json?screen_name=%s&exclude_replies=true&include_rts=false&count=50"
	urls.Twitter.Queries.BungieHelp = fmt.Sprintf(query_template, "BungieHelp")
	urls.Twitter.Queries.Destiny2Team = fmt.Sprintf(query_template, "Destiny2Team")
}
