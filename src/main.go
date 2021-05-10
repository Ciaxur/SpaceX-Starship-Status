package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/spf13/viper"
	"spacex-status.twitterapi/src/cli"
	"spacex-status.twitterapi/src/twitter"
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
func getTweets(token string, userID string) twitter.Tweet {
	// Construct Request
	var requestBody bytes.Buffer
	url := fmt.Sprintf("https://api.twitter.com/2/users/%s/tweets?tweet.fields=entities", userID)
	req, _ := http.NewRequest("GET", url, &requestBody)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Sumbit Request
	fmt.Printf("URL Request: '%s'\n", url)
	client := http.Client{}
	res, err := client.Do(req)
	handleErr(err, "Client Request Error")

	// Parse JSON Body
	var result twitter.Tweet
	json.NewDecoder(res.Body).Decode(&result)
	return result
}

func main() {
	// Init Environment
	initEnv()
	bearerToken := viper.Get("TWITTER_BEARER_TOKEN").(string)
	userID := viper.Get("USER_ID").(string)

	// Load in Cache if any
	var cache twitter.Tweet
	data, err := ioutil.ReadFile("cached.json")
	if err == nil {
		json.Unmarshal(data, &cache)
	}

	// Handle CLI Arguments
	argsForwarded := cli.HandleCliArgs(&cache)

	// Request Tweets
	tweets := getTweets(bearerToken, userID)
	latestMatch := cache.LatestMatch

	// Check if there is a New Tweet
	if tweets.Meta.NewestID != cache.Meta.NewestID {
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

				// Score Finding!
				for index := range found {
					value := string(found[index])
					switch strings.ToLower(value) {
					case "status":
						score += 0.8
					case "superheavy":
						score += 0.1
					case "spacex":
						score += 0.1
					case "starship":
						score += 0.1
					}
				}

				// Stop Loop, since Found what was looking for!
				break
			}
		}

		// Check Scoring
		fmt.Printf("Score of '%.2f'\n", score)
		if score > 0.8 && latestMatch.ID != cache.LatestMatch.ID {
			fmt.Println("New Tweet:")
			fmt.Println("- ID: ", latestMatch.ID)
			fmt.Println("- Tweet: ", latestMatch.Text)

			// Execute given method with Tweet ID and Text at end
			// Args[1] 	= Command to Execute
			// Args[2:] = Arguments to Command (Optional)
			if len(os.Args) > 1 {
				arr := append(os.Args[argsForwarded+2:], latestMatch.ID, latestMatch.Text)

				cmd := exec.Command(os.Args[argsForwarded+1], arr...)
				err := cmd.Start()
				if err != nil {
					fmt.Printf("Given Command failed to Execute: %v\n", err)
				}
			}
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
