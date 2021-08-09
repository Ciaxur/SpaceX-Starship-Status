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
	helpers "spacex-status.twitterapi/src/helpers"
	spacexdata "spacex-status.twitterapi/src/spacex-data"
	"spacex-status.twitterapi/src/twitter"
)

const (
	SCORE_TRESHOLD = 0.8
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
	helpers.HandleGeneralErr(err, "Client Request Error")

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
	argsForwarded, args := cli.HandleCliArgs(&cache)

	// Launch Check-Only
	if args.CheckLaunch {
		state := spacexdata.Init(args)
		os.Exit(state)
	}

	// Request Tweets
	tweets := getTweets(bearerToken, userID)
	latestMatch := cache.LatestMatch

	// Check if there is a New Tweet
	if tweets.Meta.NewestID != cache.Meta.NewestID {
		// Find Matching Tweet
		re, err := regexp.Compile("(?i)status|spacex|starship|production|starbase|diagram")
		helpers.HandleGeneralErr(err, "Regular Expression Failed to Compile")
		score := 0.0

		for tweetIndex := range tweets.Data {
			tweet := tweets.Data[tweetIndex]

			if found := re.FindAll([]byte(tweet.Text), -1); len(found) > 0 {
				tweet := tweets.Data[tweetIndex]
				fmt.Println("Tweet:", tweet.Text)

				// Check if already checked
				if cache.LatestMatch.ID == tweet.ID {
					score = 0.0
					break
				}

				// Keep track of Latest Match
				latestMatch = tweet

				// Score Finding!
				for _, value := range found {
					switch strings.ToLower(string(value)) {
					case "production":
						score += 0.8
					case "starbase":
						score += 0.2
					case "diagram":
						score += 0.4
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

				// Stop Loop, IF Found what was looking for!
				if score > SCORE_TRESHOLD {
					fmt.Println("Score Threashold reached")
					break
				} else {
					fmt.Println("Resetting Score to 0")
					score = 0
				}

			}
		}

		// Check Scoring
		fmt.Printf("Score of '%.2f'\n", score)
		if score > SCORE_TRESHOLD && latestMatch.ID != cache.LatestMatch.ID {
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
		helpers.HandleGeneralErr(err, "Error Converting Cache to Bytes")
		ioutil.WriteFile("cached.json", data, 0664)
	} else {
		fmt.Println("No New Tweats")
	}
}
