package main

import (
	"github.com/gtuk/discordwebhook"
)

// Sends a message via discord webhook containing the post title and link
func sendDiscordWebhook(content string) {
	username := "Bungie Alerter"
	message := discordwebhook.Message{
		Username: &username,
		Content:  &content,
	}

	// Bail if notifications are disabled
	if !notifications {
		return
	}

	// Send the notification
	err := discordwebhook.SendMessage(urls.Discord.WebhookUrl, message)
	if err != nil {
		ErrorLogger.Fatalln(err)
	} else {
		InfoLogger.Println("Fired webhook")
	}

}
