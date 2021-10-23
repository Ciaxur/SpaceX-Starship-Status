package spacexdata

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	cli "spacex-status.twitterapi/src/cli"
	"spacex-status.twitterapi/src/helpers"
)

// Downloads and returns image path
func downloadImage(url string) *os.File {
	fmt.Printf("Downloading Image %s\n", url)

	// Request Image
	client := http.Client{}
	res, err := client.Get(url)
	helpers.HandleGeneralErr(err, "Image Download Error")
	defer res.Body.Close()

	// Read Image
	rawImage, err := ioutil.ReadAll(res.Body)
	helpers.HandleGeneralErr(err, "Could not read image body")
	fmt.Printf("Image Size: %.2fKB\n", float32(len(rawImage))/(1024))

	// Save Image
	tmpFile, err := ioutil.TempFile(os.TempDir(), "spacex-*.png")
	helpers.HandleGeneralErr(err, "Could not create temp file")
	tmpFile.Write(rawImage)

	return tmpFile
}

func HandleLaunchCheck(latestLaunch *LaunchResponse, rocket *RocketResponse) {
	// Check cache if actually new launch
	cache, found := initCache()

	// No new Launch
	if found && latestLaunch.Id == cache.LaunchResponse.Id {
		fmt.Println("No new launch")
		return
	}

	// No Image associated with Launch (no launch)
	if latestLaunch.Links.Patch.Small == "" {
		fmt.Println("No image to download")
		return
	}

	imageFile := downloadImage(latestLaunch.Links.Patch.Small)
	imageFile.Close()
	fmt.Println("Image Downloaded to:", imageFile.Name())

	// Issue notification with data
	if len(os.Args) > 1 {
		flickrLinksStr := ""
		maxFlickrLinks := 2
		for idx, flickrLink := range latestLaunch.Links.Flickr.Original {
			flickrLinksStr += "\n - " + flickrLink
			if idx >= maxFlickrLinks {
				break
			}
		}

		arr := append(
			os.Args[3:],
			"-i", imageFile.Name(),
			fmt.Sprintf("%s [%s]", latestLaunch.Name, latestLaunch.Id),
			fmt.Sprintf("%s %s \n - Webcast: %s", latestLaunch.Details, flickrLinksStr, latestLaunch.Links.Webcast),
		)

		cmd := exec.Command(os.Args[2], arr...)
		err := cmd.Start()
		if err != nil {
			fmt.Printf("Given Command failed to Execute: %v\n", err)
		}
	}

	// Save Cache
	cache.LaunchResponse = latestLaunch
	cache.RocketResponse = rocket
	saveCache(*cache)
}

func Init(args *cli.Arguments) int {
	if args.CheckLaunch {
		latestLaunch_res := getLatestLaunch()
		rocket_res, err := getRocket(latestLaunch_res.RocketId)
		if err != nil {
			fmt.Println(err.Error())
			return 1
		}
		HandleLaunchCheck(&latestLaunch_res, &rocket_res)
	} else {
		fmt.Printf("SpaceX Data: No Argument Request Handled\n")
		return 1
	}

	return 0
}
