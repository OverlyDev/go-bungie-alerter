package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/joho/godotenv"
	"github.com/mmcdole/gofeed"
	"github.com/urfave/cli/v2"
)

func censorWebhook(webhook string) string {
	lastIndex := strings.LastIndex(webhook, "/")
	public := webhook[:lastIndex]
	private := webhook[lastIndex+1:]
	return public + "/****" + private[len(private)-4:]
}

func cliHandler() {
	app := &cli.App{
		Name:    "BungieAlerter",
		Version: fmt.Sprintf("%s | ref: %s | built: %s\n\t%s", version, reference, buildTime, "Repo: https://github.com/OverlyDev/go-bungie-alerter"),
		Usage:   "Sends messages to Discord webhooks on new Bungie posts",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "webhook",
				Aliases: []string{"w"},
				Usage:   "discord webhook `URL` where notifications will be sent",
				EnvVars: []string{"DISCORD_WEBHOOK"},
			},
			&cli.BoolFlag{
				Name:    "silent",
				Aliases: []string{"s"},
				Usage:   "disables webhook usage",
			},
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "logs additional information",
			},
			&cli.BoolFlag{
				Name:    "logfile",
				Aliases: []string{"l"},
				Usage:   "Enable logging to file",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "go",
				Usage: "Start BungieAlerter",
				Action: func(cCtx *cli.Context) error {
					printVersion()
					handleFlags(cCtx)

					startTime = time.Now().UTC()
					readTimestampsFile()
					actionLoop()
					return nil
				},
			},
			{
				Name:  "update",
				Usage: "Check for new version of BungieAlerter",
				Action: func(cCtx *cli.Context) error {
					// handleFlags(cCtx)
					checkForNewVersion()
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func obtainWebhookUrl(cCtx *cli.Context) {
	// Obtain webhook url from flags, falling back to env
	webhook := cCtx.String("webhook")
	source := ""

	if webhook != "" {
		source = "flag"
	} else {
		err := godotenv.Load()
		if err != nil {
			ErrorLogger.Println("Error loading env")
			os.Exit(1)
		}
		_, found := os.LookupEnv("DISCORD_WEBHOOK")
		if !found {
			ErrorLogger.Println("DISCORD_WEBHOOK not provided, exiting")
			os.Exit(1)
		} else {
			webhook = os.Getenv("DISCORD_WEBHOOK")
			source = "env"
		}
	}

	// Validate webhook, exit if invalid
	if !govalidator.IsURL(webhook) {
		ErrorLogger.Println("Invalid webhook:", webhook)
		os.Exit(1)
	}

	// Censor the webhook for console print
	censored := censorWebhook(webhook)
	InfoLogger.Printf("Got webhook from %s: %s\n", source, censored)

	// Save the webhook for later use
	urls.Discord.WebhookUrl = webhook
}

func actionLoop() {
	feedParser := gofeed.NewParser()

	InfoLogger.Println("Starting")
	for {
		InfoLogger.Println("Getting Bungie.net feed")
		newPost := parseBungiePosts(feedParser)

		InfoLogger.Println("Getting tweets")
		newTweet := checkForTweets()

		if newPost || newTweet {
			InfoLogger.Println("Changes to timestamps, writing to disk")
			writeTimestampsFile()
		} else {
			InfoLogger.Println("No changes to timestamps")
		}
		InfoLogger.Println("Sleeping 60s")
		time.Sleep(60 * time.Second)
		fmt.Println()
	}
}

func handleFlags(context *cli.Context) {
	// Enable debug logging if given debug flag
	if context.Bool("debug") {
		setLoggerDebug()
		DebugLogger.Println("Debug logs enabled")
	}

	// Disable webhooks if given silent flag
	if !context.Bool("silent") {
		notifications = true
		obtainWebhookUrl(context)
	} else {
		notifications = false
		InfoLogger.Println("Webhook notifications disabled")
	}

	// Enable logging to file if given logfile flag
	if context.Bool("logfile") {
		setLoggerMultiWrite(context.Bool("debug"))
	}
}
