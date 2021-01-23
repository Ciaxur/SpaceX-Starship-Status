package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

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
	latestMatch := cache.LatestMatch

	// Check if there is a New Tweet
	if tweets.Meta.NewestID != cache.Meta.NewestID {
		// DEBUG: Logs
		fmt.Println("New Tweet: ", tweets.Meta.NewestID)
		fmt.Println("Tweet ID: ", tweets.Data[0].ID)

		// Find Matching Tweet
		re, err := regexp.Compile("(?i)status|spacex|starship")
		handleErr(err, "Regular Expression Failed to Compile")
		score := 0.0

		for tweetIndex := range tweets.Data {
			tweet := tweets.Data[tweetIndex]

			if found := re.FindAll([]byte(tweet.Text), -1); len(found) > 0 {
				tweet := tweets.Data[tweetIndex]

				// Check if already checked
				if cache.LatestMatch.ID == tweet.ID {
					score = 0.0
					break
				}

				// Keep track of Latest Match
				latestMatch = tweet
				fmt.Printf("%s: %s\n", tweet.ID, tweet.Text)
				fmt.Printf("Found: %s\n\n", found)

				// Score Finding!
				for index := range found {
					value := string(found[index])
					switch strings.ToLower(value) {
					case "status":
						score += 0.8
						break
					case "spacex":
						score += 0.1
						break
					case "starship":
						score += 0.1
						break
					}
				}

				// Stop Loop, since Found what was looking for!
				break
			}
		}

		// Check Scoring
		fmt.Printf("Score of '%.2f'\n", score)
		if score > 0.8 && latestMatch.ID != cache.LatestMatch.ID {
			fmt.Println("Latest Matched ID: ", latestMatch.ID)
		} else {
			latestMatch = cache.LatestMatch // Keep older Match
		}

		// Store current New State
		cache = tweets
		cache.LatestMatch = latestMatch
		data, err := json.Marshal(cache)
		handleErr(err, "Error Converting Cache to Bytes")
		ioutil.WriteFile("cached.json", data, 0664)
	} else {
		fmt.Println("No New Tweats")
	}

}
