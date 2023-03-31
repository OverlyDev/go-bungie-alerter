package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

// Thanks, internet
var publicToken = "Bearer AAAAAAAAAAAAAAAAAAAAAPYXBAAAAAAACLXUNDekMxqa8h%2F40K4moUkGsoc%3DTYfbDKbT3jJPCEVnMYqilB28NHfOPqkca3qaAxGfsyKCs0wRbw"

type guestTokenStruct struct {
	GuestToken string `json:"guest_token"`
}

type twitterAuthStruct struct {
	Public string
	Guest  string
}

type tweetStruct struct {
	Created string `json:"created_at"`
	Id      string `json:"id_str"`
	Text    string `json:"text"`
}

// Obtains an auth token from a public token
func getTwitterAuth() {
	// Create the request
	req, err := http.NewRequest("POST", urls.Twitter.Auth, nil)
	if err != nil {
		ErrorLogger.Println(err)
	}

	// Add our header with the public token
	req.Header.Set("Authorization", publicToken)

	// Send request
	DebugLogger.Println("Sending twitter auth request")
	client := http.DefaultClient
	response, err := client.Do(req)
	if err != nil {
		ErrorLogger.Println(err)
		DebugLogger.Println("Response:", response)
	}

	// Read the body
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		ErrorLogger.Println(err)
		DebugLogger.Println("ResponseData:", responseData)
	}

	// Unmarshal json response
	var responseObject guestTokenStruct
	json.Unmarshal(responseData, &responseObject)

	// Save token
	twitterAuth.Public = publicToken
	twitterAuth.Guest = responseObject.GuestToken
	DebugLogger.Println("Twitter Token:", twitterAuth.Guest)

}

// Returns latest tweet from supplied url
func makeTwitterRequest(url, account string) (*tweetStruct, error) {
	// Create the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		ErrorLogger.Println(err)
		return nil, err
	}

	// Add our auth headers, acquired via getTwitterAuth
	req.Header.Set("Authorization", twitterAuth.Public)
	req.Header.Set("x-guest-token", twitterAuth.Guest)
	DebugLogger.Println("Headers:", req.Header)

	//Send request
	DebugLogger.Println("Sending twitter tweet request")
	client := http.DefaultClient
	response, err := client.Do(req)
	if err != nil {
		ErrorLogger.Println(err)
		DebugLogger.Println("Response:", response)
		return nil, err
	}

	// Something went wrong :(
	if response.StatusCode != 200 {
		ErrorLogger.Printf("Got bad twitter response from: %s\n", account)
		DebugLogger.Println("Response:", response)
		return nil, fmt.Errorf("no data")
	}

	// Read the body
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		ErrorLogger.Println(err)
		DebugLogger.Println("ResponseData:", responseData)
		return nil, err
	}

	// Unmarshall json response
	var responseObject []tweetStruct
	json.Unmarshal(responseData, &responseObject)

	// Grab and return the latest tweet
	latestTweet := responseObject[0]
	DebugLogger.Println("LatestTweet:", latestTweet)
	return &latestTweet, nil
}

// Checks for new tweets using all queries stored in urls.Twitter.Queries
func checkForTweets() bool {
	changes := false
	v := reflect.ValueOf(urls.Twitter.Queries)
	vType := v.Type()

	// Iterate through urls.Twitter.Queries
	for i := 0; i < v.NumField(); i++ {
		account := vType.Field(i).Name
		url := v.Field(i).String()

		// Failed to get the tweet :(
		tweet, err := makeTwitterRequest(url, account)
		if err != nil {
			ErrorLogger.Printf("Error getting tweets for: %s\n", account)
			DebugLogger.Println("Tweet:", tweet)
			continue
		}

		tweetTime := convertTwitterTimeStrToTime(tweet.Created)
		DebugLogger.Println("TweetTime:", tweetTime)
		lastTweetTime := convertStrToTime(getField(&timestamps, "Twitter"+account))
		DebugLogger.Println("LastTweetTime:", lastTweetTime)

		// No new tweets
		if tweetTime.Before(lastTweetTime) || tweetTime.Equal(lastTweetTime) {
			InfoLogger.Printf("Up to date: %s\n", account)
			changes = changes || false

			// There do be a new tweet
		} else {
			AlertLogger.Printf("New tweet from: %s\n", account)
			content := fmt.Sprintf("New tweet from %s\n", account)
			content += fmt.Sprintf(urls.Twitter.TweetTemplate, account, tweet.Id)
			DebugLogger.Println("Content:", content)
			sendDiscordWebhook(content)
			newTimestamp := convertTimeToStr(tweetTime)
			DebugLogger.Println("NewTimestamp:", newTimestamp)

			switch account {
			case "BungieHelp":
				timestamps.TwitterBungieHelp = newTimestamp
				DebugLogger.Println("Updated timestamps.TwitterBungieHelp with new timestamp")
			case "Destiny2Team":
				timestamps.TwitterDestiny2Team = newTimestamp
				DebugLogger.Println("Updated timestamps.TwitterDestiny2Team with new timestamp")
			}
			changes = changes || true
			DebugLogger.Println("Changes:", changes)
		}
	}
	return changes

}
