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
	jsonErr := json.Unmarshal(responseData, &responseObject)
	if err != nil {
		ErrorLogger.Println("Error during json.Unmarshal")
		DebugLogger.Println(jsonErr)
	}

	// Save token
	twitterAuth.Public = publicToken
	twitterAuth.Guest = responseObject.GuestToken
	DebugLogger.Println("Twitter Token:", twitterAuth.Guest)

}

// Creates a request object containing given url and our headers
//
//	success = (http.Request, nil)
//	fail = (nil, err)
func createTwitterWebRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		ErrorLogger.Println("Failed to create Twitter web request for:", url)
		DebugLogger.Println(err)
		return nil, err
	}

	req.Header.Set("Authorization", twitterAuth.Public)
	req.Header.Set("x-guest-token", twitterAuth.Guest)

	return req, nil
}

// "Does" a request
//
//	success = (reponse.statusCode, response.Body)
//	fail = (0, nil)
func submitWebRequest(request *http.Request) (int, io.ReadCloser) {
	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		ErrorLogger.Println("Error making web request to:", request.URL)
		DebugLogger.Println(err)
		return 0, nil
	}

	return response.StatusCode, response.Body
}

// Returns latest tweet from supplied url
func makeTwitterRequest(url, account string) (*tweetStruct, error) {
	var statusCode int = 0
	var body io.ReadCloser = nil

	// Handle non-200 http status codes
	attempts := 1
	for statusCode != 200 {
		// Bail if more than 3 attempts
		if attempts > 3 {
			ErrorLogger.Println("Max retry attempts")
			return nil, fmt.Errorf("no data")
		}

		// Create the request
		DebugLogger.Println("Creating request for:", url)
		req, err := createTwitterWebRequest("GET", url)
		if err != nil {
			return nil, fmt.Errorf("no data")
		}
		// Send request
		DebugLogger.Println("Sending twitter tweet request")
		statusCode, body = submitWebRequest(req)

		switch statusCode {
		// Bail on a 404
		case 404:
			ErrorLogger.Println("Got 404 from account:", account)
			DebugLogger.Println("Full url:", url)
			return nil, fmt.Errorf("no data")

		// Renew guest-token on a 403
		case 403:
			InfoLogger.Println("Renewing Twitter guest-token")
			getTwitterAuth()
		}

		attempts += 1
	}

	// Read the body
	responseData, err := io.ReadAll(body)
	if err != nil {
		ErrorLogger.Println(err)
		DebugLogger.Println("ResponseData:", responseData)
		return nil, err
	}

	// Unmarshall json response
	var responseObject []tweetStruct
	jsonErr := json.Unmarshal(responseData, &responseObject)
	if jsonErr != nil {
		ErrorLogger.Println("Error during json.Unmarshal")
		DebugLogger.Println(jsonErr)
	}

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

		// Compare timestamps
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
