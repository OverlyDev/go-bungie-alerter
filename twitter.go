package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

// Thanks, internet
var public_token = "Bearer AAAAAAAAAAAAAAAAAAAAAPYXBAAAAAAACLXUNDekMxqa8h%2F40K4moUkGsoc%3DTYfbDKbT3jJPCEVnMYqilB28NHfOPqkca3qaAxGfsyKCs0wRbw"

type guest_token_struct struct {
	Guest_token string `json:"guest_token"`
}

type twitter_auth_struct struct {
	Public string
	Guest  string
}

type tweet_struct struct {
	Created  string `json:"created_at"`
	Id       string `json:"id_str"`
	Text     string `json:"text"`
	Entities struct {
		Urls []struct {
			ExpandedUrl string `json:"expanded_url"`
		} `json:"urls"`
	} `json:"entities"`
}

// Obtains an auth token from a public token
func get_twitter_auth() {
	// Create the request
	req, err := http.NewRequest("POST", urls.Twitter.Auth, nil)
	if err != nil {
		ErrorLogger.Println(err)
	}

	// Add our header with the public token
	req.Header.Set("Authorization", public_token)

	// Send request
	client := http.DefaultClient
	response, err := client.Do(req)
	if err != nil {
		ErrorLogger.Println(err)
	}

	// Read the body
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		ErrorLogger.Println(err)
	}

	// Unmarshal json response
	var responseObject guest_token_struct
	json.Unmarshal(responseData, &responseObject)

	// Save token
	twitter_auth.Public = public_token
	twitter_auth.Guest = responseObject.Guest_token

}

// Returns latest tweet from supplied url
func make_twitter_request(url, account string) (*tweet_struct, error) {
	// Create the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		ErrorLogger.Println(err)
		return nil, err
	}

	// Add our auth headers, acquired via get_twitter_auth
	req.Header.Set("Authorization", twitter_auth.Public)
	req.Header.Set("x-guest-token", twitter_auth.Guest)

	//Send request
	client := http.DefaultClient
	response, err := client.Do(req)
	if err != nil {
		ErrorLogger.Println(err)
		return nil, err
	}

	// Something went wrong :(
	if response.StatusCode != 200 {
		ErrorLogger.Printf("Got bad twitter response from: %s\n", account)
		return nil, fmt.Errorf("no data")
	}

	// Read the body
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		ErrorLogger.Println(err)
		return nil, err
	}

	// Unmarshall json response
	var responseObject []tweet_struct
	json.Unmarshal(responseData, &responseObject)

	// Grab and return the latest tweet
	latest_tweet := responseObject[0]
	return &latest_tweet, nil
}

// Checks for new tweets using all queries stored in urls.Twitter.Queries
func check_for_tweets() bool {
	changes := false
	v := reflect.ValueOf(urls.Twitter.Queries)
	vType := v.Type()

	// Iterate through urls.Twitter.Queries
	for i := 0; i < v.NumField(); i++ {
		account := vType.Field(i).Name
		url := v.Field(i).String()

		// Failed to get the tweet :(
		tweet, err := make_twitter_request(url, account)
		if err != nil {
			continue
		}

		tweet_time := convert_twitter_time_str_to_time(tweet.Created)
		last_tweet_time := convert_RFC1123_str_to_time(timestamps.TwitterBungieHelp)

		// No new tweets
		if tweet_time.Before(last_tweet_time) || tweet_time.Equal(last_tweet_time) {
			InfoLogger.Printf("Up to date: %s\n", account)
			changes = changes || false

			// There do be a new tweet
		} else {
			AlertLogger.Printf("New tweet from: %s\n", account)
			content := fmt.Sprintf("New tweet from %s\n", account)
			content += tweet.Entities.Urls[0].ExpandedUrl
			send_discord_webhook(content)
			new_timestamp := convert_time_to_RFC1123_str(tweet_time)
			switch account {
			case "BungieHelp":
				timestamps.TwitterBungieHelp = new_timestamp
			case "Destiny2Team":
				timestamps.TwitterDestiny2Team = new_timestamp
			}
			changes = changes || true
		}
	}
	return changes

}
