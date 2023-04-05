package main

import (
	"io"
	"log"
	"time"
)

var (
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	AlertLogger   *log.Logger
	DebugLogger   *log.Logger
	urls          urlStorageStruct
	twitterAuth   twitterAuthStruct
	timestamps    timestampStruct
	notifications bool
	startTime     time.Time
)

func init() {
	// Set up graceful exits
	signalHandler()

	// Set up the loggers
	InfoLogger = log.New(log.Default().Writer(), "INFO | ", log.Ldate|log.Ltime|log.LUTC)
	ErrorLogger = log.New(log.Default().Writer(), "ERRO | ", log.Ldate|log.Ltime|log.LUTC)
	AlertLogger = log.New(log.Default().Writer(), "ALRT | ", log.Ldate|log.Ltime|log.LUTC)
	DebugLogger = log.New(io.Discard, "DBUG | ", log.Default().Flags())

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
