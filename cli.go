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
		Name: "BungieAlerter",
		Version: fmt.Sprintf("%s | commit: %s | built: %s",
			strings.ReplaceAll(version, "\n", ""),
			strings.ReplaceAll(commit, "\n", ""),
			strings.ReplaceAll(buildTime, "\n", ""),
		),
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "OverlyDev",
				Email: "kyle@overly.dev",
			},
		},
		Usage: "Sends messages to Discord webhooks on new Bungie posts",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "webhook",
				Aliases: []string{"w"},
				Usage:   "Discord webhook `URL` where notifications will be sent",
				EnvVars: []string{"DISCORD_WEBHOOK"},
			},
			&cli.BoolFlag{
				Name:    "silent",
				Aliases: []string{"s"},
				Usage:   "Disables webhook usage",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "go",
				Usage: "Start BungieAlerter",
				Action: func(cCtx *cli.Context) error {
					printVersion()
					obtainWebhookUrl(cCtx)
					readTimestampsFile()
					action_loop()
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func obtainWebhookUrl(cCtx *cli.Context) error {
	// Disable webhooks if given silent flag
	// Also bypass obtaining webhook url since it won't be used
	if cCtx.Bool("silent") {
		notifications = false
		InfoLogger.Println("Webhook notifications disabled")
		return nil
	} else {
		notifications = true
	}

	// Obtain webhook url from flags, falling back to env
	webhook := cCtx.String("webhook")
	source := ""

	if webhook != "" {
		source = "flag"
	} else {
		godotenv.Load()
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
	return nil
}

func action_loop() {
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
	}
}
