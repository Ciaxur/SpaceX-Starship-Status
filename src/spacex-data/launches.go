package spacexdata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	helpers "spacex-status.twitterapi/src/helpers"
)

type LaunchResponse struct {
	Fairings struct {
		Reused          bool `json:"reused"`
		RecoveryAttempt bool `json:"recovery_attempt"`
		Recovered       bool `json:"recovered"`
	} `json:"fairings"`

	Links struct {
		Patch struct {
			Small string `json:"small"`
			Largs string `json:"large"`
		} `json:"patch"`

		Reddit struct {
			Campaign string `json:"campaign"`
			Launch   string `json:"launch"`
			Media    string `json:"media"`
			Recovery string `json:"recovery"`
		} `json:"reddit"`

		Flickr struct {
			Small    []string `json:"small"`
			Original []string `json:"original"`
		} `json:"flickr"`

		Webcast string `json:"webcast"`
	} `json:"links"`

	StaticFireDateUTC  string `json:"static_fire_date_utc"`
	StaticFireDateUnix int64  `json:"static_fire_date_unix"`
	TBD                bool   `json:"tbd"`
	Net                bool   `json:"net"`
	RocketId           string `json:"rocket"`
	Success            bool   `json:"success"`
	Details            string `json:"details"`
	FlightNumber       int64  `json:"flight_number"`
	Name               string `json:"name"`
	DateUTC            string `json:"date_utc"`
	DateUnix           string `json:"date_unix"`
	Upcoming           bool   `json:"upcoming"`
	Id                 string `json:"id"`
}

// Fetches Latest Launch Data
func getLatestLaunch() LaunchResponse {
	url := "https://api.spacexdata.com/v5/launches/latest"
	var requestBody bytes.Buffer
	req, _ := http.NewRequest("GET", url, &requestBody)

	// Submit Request
	fmt.Printf("Requesting Latest Launch: %s\n", url)
	client := http.Client{}
	res, err := client.Do(req)
	helpers.HandleGeneralErr(err, "Latest Launch Request Error")

	// Parse JSON Body
	var result LaunchResponse
	json.NewDecoder(res.Body).Decode((&result))
	return result
}
