package spacexdata

import (
	"encoding/json"
	"io/ioutil"

	"spacex-status.twitterapi/src/helpers"
)

type SpaceXDataCache struct {
	LaunchResponse *LaunchResponse `json:"launchResponse"`
	RocketResponse *RocketResponse `json:"rocketResponse"`
}

var (
	spacexDataCache SpaceXDataCache
	cacheFilename   string = "spacex-data-cache.json"
)

// Instantiates Cache, returning Cache Struct & state of cache found
func initCache() (*SpaceXDataCache, bool) {
	// Load Cache IF available
	data, err := ioutil.ReadFile(cacheFilename)
	if err == nil {
		json.Unmarshal(data, &spacexDataCache)
		return &spacexDataCache, true
	}
	return &spacexDataCache, false
}

// Saves given Cache Struct to file
func saveCache(cache SpaceXDataCache) {
	data, err := json.Marshal(cache)
	helpers.HandleGeneralErr(err, "Error Marshalling SpaceX Data Cache")
	ioutil.WriteFile(cacheFilename, data, 0644)
}
