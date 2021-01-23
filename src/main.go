package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

/**
 * Initializes environment variables from
 *  the .env file
 */
func initEnv() {
	// Load in Environmental Configuration
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error reading in '.env' file: %s", err))
	}
}

// Handle Errors
func handleErr(err error, outStr string) {
	if err != nil {
		panic(fmt.Errorf("%s: %s", outStr, err))
	}
}

/**
 * Requests a tweet data from API
 * @param token Bearer Token for Twitter's API
 * @param userID The user's ID to request Tweets of
 */
func getTweets(token string, userID string) Tweet {
	// Construct Request
	var requestBody bytes.Buffer
	url := fmt.Sprintf("https://api.twitter.com/2/users/%s/tweets", userID)
	req, _ := http.NewRequest("GET", url, &requestBody)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Sumbit Request
	fmt.Printf("URL Request: '%s'\n", url)
	client := http.Client{}
	res, err := client.Do(req)
	handleErr(err, "Client Request Error")

	// Parse JSON Body
	var result Tweet
	json.NewDecoder(res.Body).Decode(&result)
	return result
}

func main() {
	// Init Environment
	initEnv()
	bearerToken := viper.Get("TWITTER_BEARER_TOKEN").(string)
	userID := viper.Get("USER_ID").(string)

	// Load in Cache if any
	var cache Tweet
	data, err := ioutil.ReadFile("cached.json")
	if err == nil {
		json.Unmarshal(data, &cache)
	}

	// Request Tweets
	tweets := getTweets(bearerToken, userID)

	// Check if there is a New Tweet
	if tweets.Meta.NewestID != cache.Meta.NewestID {
		log.Println("New Tweet: ", tweets.Meta.NewestID)
		log.Println("Tweet ID: ", tweets.Data[0].ID)
		log.Println("Tweet: ", tweets.Data[0].Text)

		// Save Cache
		cache = tweets
		data, err := json.Marshal(cache)
		handleErr(err, "Error Converting Cache to Bytes")
		ioutil.WriteFile("cached.json", data, 0664)
	}

}
