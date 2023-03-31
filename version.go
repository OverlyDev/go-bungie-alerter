package main

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:generate bash gen_embed_files.sh
var (
	//go:embed embeds/time.txt
	buildTime string
	//go:embed embeds/version.txt
	version string
	//go:embed embeds/commit.txt
	commit string
)

func printVersion() {
	fmt.Printf(
		"BungieAlerter | Version: %s | Commit: %s | Built: %s\n",
		strings.ReplaceAll(version, "\n", ""),
		strings.ReplaceAll(commit, "\n", ""),
		strings.ReplaceAll(buildTime, "\n", ""),
	)
}
