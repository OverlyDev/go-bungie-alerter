package main

import (
	"github.com/gtuk/discordwebhook"
)

// Sends a message via discord webhook containing the post title and link
func sendDiscordWebhook(content string) {
	// Bail if notifications are disabled
	if !notifications {
		return
	}

	username := "Bungie Alerter"
	message := discordwebhook.Message{
		Username: &username,
		Content:  &content,
	}
	DebugLogger.Println("MessageUsername:", *message.Username)
	DebugLogger.Println("MessageContent:", *message.Content)

	// Send the notification
	err := discordwebhook.SendMessage(urls.Discord.WebhookUrl, message)
	if err != nil {
		ErrorLogger.Fatalln(err)
	} else {
		InfoLogger.Println("Sent webhook notification")
	}

}
