package main

import (
	"log"
)

var (
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	AlertLogger   *log.Logger
	urls          urlStorageStruct
	twitterAuth   twitterAuthStruct
	timestamps    timestampStruct
	notifications bool
)

func init() {
	// Set up graceful exits
	signalHandler()

	// Set up the loggers
	InfoLogger = log.New(log.Default().Writer(), "INFO | ", log.Ldate|log.Ltime|log.Lshortfile|log.LUTC)
	ErrorLogger = log.New(log.Default().Writer(), "ERROR | ", log.Ldate|log.Ltime|log.Lshortfile|log.LUTC)
	AlertLogger = log.New(log.Default().Writer(), "ALERT | ", log.Ldate|log.Ltime|log.Lshortfile|log.LUTC)

	// Set up all the urls
	populateUrlStorage()

	// Twitter auth setup stuff
	getTwitterAuth()
	if twitterAuth.Guest == "" {
		ErrorLogger.Fatalln("Failed to obtain twitter auth token, exiting")
	}
}

func main() {
	cliHandler()
}
