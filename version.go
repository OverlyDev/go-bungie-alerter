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
	//go:embed embeds/ref.txt
	reference string
)

func printVersion() {
	fmt.Printf(
		"BungieAlerter | Version: %s | Ref: %s | Built: %s\n",
		strings.ReplaceAll(version, "\n", ""),
		strings.ReplaceAll(reference, "\n", ""),
		strings.ReplaceAll(buildTime, "\n", ""),
	)
}
