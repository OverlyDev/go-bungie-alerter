package main

import (
	"github.com/gtuk/discordwebhook"
)

// Sends a message via discord webhook containing the post title and link
func send_discord_webhook(content string) {
	username := "BotUser"

	message := discordwebhook.Message{
		Username: &username,
		Content:  &content,
	}

	err := discordwebhook.SendMessage(urls.Discord.WebhookUrl, message)
	if err != nil {
		ErrorLogger.Fatalln(err)
	} else {
		InfoLogger.Println("Fired webhook")
	}

}
