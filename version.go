package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/mod/semver"
)

//go:generate bash gen_embed_files.sh
var (
	//go:embed embeds/time.txt
	buildTime string
	//go:embed embeds/version.txt
	version string
	//go:embed embeds/ref.txt
	reference string
)

type latestGithubVersion struct {
	Url       string `json:"html_url"`
	Tag       string `json:"tag_name"`
	Published string `json:"published_at"`
}

// Checks github for a newer version of BungieAlerter
func checkForNewVersion() {
	req, err := http.NewRequest("GET", "https://api.github.com/repos/OverlyDev/go-bungie-alerter/releases/latest", nil)
	if err != nil {
		ErrorLogger.Println("Aborting version check")
		return
	}

	status, body := submitWebRequest(req)
	if status != 200 {
		ErrorLogger.Println("Failed to get latest version info, aborting version check")
		return
	}

	bodyData, err := io.ReadAll(body)
	if err != nil {
		ErrorLogger.Println("Error reading body, aborting version check")
		DebugLogger.Println("body:", body)
		DebugLogger.Println("bodyData:", bodyData)
		return
	}

	var bodyJson latestGithubVersion
	jsonErr := json.Unmarshal(bodyData, &bodyJson)
	if jsonErr != nil {
		ErrorLogger.Println("Error during json.Unmarshal, aborting version check")
		DebugLogger.Println(jsonErr)
	}

	printVersion()

	switch semver.Compare(version, bodyJson.Tag) {
	case 0:
		fmt.Println("\nYou're on the latest version! :)")
	case -1:
		fmt.Printf("\n[!] New version [%s] available as of %s\n", bodyJson.Tag, bodyJson.Published)
		fmt.Printf("[!] Link: %s\n", bodyJson.Url)
	case 1:
		fmt.Println("You're somehow on a newer version that what's available. :huh:")
	default:
		fmt.Println("I have no clue, this should never happen. Please file a bug report :)")
	}

}

func printVersion() {
	fmt.Printf("BungieAlerter | Version: %s | Ref: %s | Built: %s\n", version, reference, buildTime)
}
