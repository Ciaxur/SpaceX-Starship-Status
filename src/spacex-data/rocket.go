package spacexdata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	helpers "spacex-status.twitterapi/src/helpers"
)

type distanceMetric struct {
	Meters float32 `json:"meters"`
	Feet   float32 `json:"feet"`
}
type weightMetric struct {
	Kg float32 `json:"kg"`
	Lb float32 `json:"lb"`
}
type thrustMetric struct {
	KiloNewton int64 `json:"kN"`
	PoundForce int64 `json:"lbf"`
}

type RocketResponse struct {
	Height   distanceMetric `json:"height"`
	Diameter distanceMetric `json:"diameter"`
	Mass     weightMetric   `json:"mass"`

	Engines struct {
		Isp struct {
			SeaLevel int32 `json:"sea_level"`
			Vacuum   int32 `json:"vacuum"`
		} `json:"isp"`

		ThrustSeaLevel thrustMetric `json:"thrust_sea_level"`
		ThrustVacuum   thrustMetric `json:"thrust_vacuum"`
		Number         int32        `json:"number"`
		Type           string       `json:"type"`
		Version        string       `json:"version"`
		Layout         string       `json:"layout"`
		Propellant1    string       `json:"propellant_1"`
		Propellant2    string       `json:"propellant_2"`
		ThrustToWeight int32        `json:"thrust_to_weight"`
	} `json:"engines"`

	FlickrImages  []string `json:"flickr_images"`
	Name          string   `json:"name"`
	Type          string   `json:"type"`
	Active        bool     `json:"active"`
	Stages        int32    `json:"stages"`
	Boosters      int32    `json:"boosters"`
	CostPerLaunch int64    `json:"cost_per_launch"`
	SuccessRate   int32    `json:"success_rate_pct"`
	Description   string   `json:"description"`
	Id            string   `json:"id"`
}

type RocketQueryResponse struct {
	Docs []RocketResponse `json:"docs"`
}

// Fetches given Rocket ID
func getRocket(id string) (RocketResponse, error) {
	url := "https://api.spacexdata.com/v4/rockets/query"
	requestBody := []byte(fmt.Sprintf(`{
		"query": {
			"_id": "%s"
		},
		"options": {}
	}`, id))
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Submit Request
	fmt.Printf("Requesting Rocket[%s]: %s\n", id, url)
	client := http.Client{}
	res, err := client.Do(req)
	helpers.HandleGeneralErr(err, "Rocket Request Error")

	// Parse JSON Body
	var result RocketQueryResponse
	json.NewDecoder(res.Body).Decode((&result))

	if len(result.Docs) == 0 {
		return RocketResponse{}, fmt.Errorf("rocket[%s] not found", id)
	}
	return result.Docs[0], nil
}
