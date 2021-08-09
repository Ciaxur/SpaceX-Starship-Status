package cli

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/fatih/color"
	"spacex-status.twitterapi/src/twitter"
)

type Arguments struct {
	IsHelp      bool
	IsVersion   bool
	IsListCache bool
	IgnoreCache bool
	CheckLaunch bool
}

// Parses CLI Arguments
//  Returning Arguments Struct of Values
func parseInput() Arguments {
	var flagHelp = flag.Bool("help", false, "Displays Help Menu")
	flag.BoolVar(flagHelp, "h", false, "Displays Help Menu")

	var flagVersion = flag.Bool("version", false, "Displays App Version")
	flag.BoolVar(flagVersion, "v", false, "Displays App Version")

	var flagList = flag.Bool("list", false, "List Cached Tweets")
	flag.BoolVar(flagList, "l", false, "List Cached Tweets")

	var flagIgnoreCache = flag.Bool("no-cache", false, "Ignores existing Cache")
	flag.BoolVar(flagIgnoreCache, "c", false, "Ignores existing Cache")

	var checkLaunch = flag.Bool("check-launch", false, "Checks latest launch information")

	flag.Parse()
	return Arguments{*flagHelp, *flagVersion, *flagList, *flagIgnoreCache, *checkLaunch}
}

// Prints Help Menu
func printHelp() {
	cyan := color.New(color.FgHiCyan).SprintFunc()

	InfoOut.Print("Usage SpaceX Starship Status:\n\t")
	fmt.Printf(cyan("SpaceX-SN-Status") + " [OPTIONS] - Initiate a Listen for SpaceX Starship Status\n")

	InfoOut.Printf("Help Options:\n")
	fmt.Printf("\t-h, -help \t\t\t Displays Help Menu\n")
	fmt.Printf("\t-v, -version \t\t\t Displays Version\n")

	InfoOut.Printf("\nCommand-Line Options:\n")
	fmt.Printf("\t-l, -list \t\t\t Lists Cached Tweets\n")
	fmt.Printf("\t-c, -no-cache \t\t\t Ignores existing Cache\n")
	fmt.Printf("\t-check-launch \t\t\t Checks latest launch information\n")
}

// Prints Tweet in a Formated Style
func printTweet(tweetStr string) {
	// Regex Patterns
	reLinks := regexp.MustCompile(`(http(s?)\:[^\s]*)`)
	reSpaceX := regexp.MustCompile(`(?i)(spacex|starship|superheavy)`)
	reEmpty := regexp.MustCompile(`(^$|\n)`)

	// Print with Style 😎
	outputResult := reLinks.ReplaceAll(
		[]byte(tweetStr),
		[]byte(InfoOut.Sprint("${1}")),
	)

	outputResult = reSpaceX.ReplaceAll(
		outputResult,
		[]byte(ErrOut.Sprint("${1}")),
	)

	outputResult = reEmpty.ReplaceAll(
		outputResult,
		[]byte(" "),
	)

	StdOut.Println(string(outputResult))
}

// HandleCliArgs -
//  Parses CLI Arguments and Handles appropriately
//  Returns the number of Arguments to ignore (It was handled here)
func HandleCliArgs(cache *twitter.Tweet) (int, *Arguments) {
	args := parseInput()

	// Handle Argument
	if args.IsHelp {
		printHelp()
		os.Exit(0)
	} else if args.IsVersion {
		StdOut.Printf("Version: %s\n", AppVersion)
		os.Exit(0)
	} else if args.IsListCache {
		if len(cache.Data) != 0 {
			// Output Latest Match
			titleOut := color.New(color.Bold)
			titleOut.Printf("Latest Match (%s)\n", cache.LatestMatch.ID)
			printTweet(cache.LatestMatch.Text)

			// Output rest of the Cache
			titleOut.Println("\nCached Tweets:")
			for index, tweet := range cache.Data {
				titleOut.Printf("[%d](%s) ", index, tweet.ID)
				printTweet(tweet.Text)
			}
		} else {
			WarnOut.Println("Cannot issue List. Cache is Empty")
		}
		os.Exit(0)
	} else if args.IgnoreCache {
		fmt.Println("Cache Ignored...")
		*cache = twitter.Tweet{}
		return 1, nil
	}

	// Nothing Handled
	return 0, &args
}
