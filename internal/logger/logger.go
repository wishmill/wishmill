package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var DebugLogger *log.Logger
var InfoLogger *log.Logger
var WarningLogger *log.Logger
var ErrorLogger *log.Logger
var FatalLogger *log.Logger

func Init() {
	//Initialize all loggers
	DebugLogger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	FatalLogger = log.New(os.Stdout, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Init2(loglevel string, devmode bool) {
	//Use loglevel debug
	if strings.ToLower(loglevel) == "debug" {
		fmt.Println("Using loglevel debug")
		return
	}

	//Disable debug logger on loglevel info
	if strings.ToLower(loglevel) == "info" {
		DebugLogger.SetOutput(io.Discard)
		fmt.Println("Using loglevel info")
		return
	}

	//Disable debug and info logger on debug level warning
	if strings.ToLower(loglevel) == "warning" {
		DebugLogger.SetOutput(io.Discard)
		InfoLogger.SetOutput(io.Discard)
		fmt.Println("Using loglevel warning")
		return
	}

	//If loglevel is not empty string now, it's not legal. Return a warning and continue with setting default log level.
	if strings.ToLower(loglevel) != "" {
		WarningLogger.Println("logger: Provided loglevel is not a legal loglevel. Using DEBUG instead.")
	}

	if devmode {
		fmt.Println("Using default loglevel for dev mode : debug")
		return
	}

	if !devmode {
		DebugLogger.SetOutput(io.Discard)
		fmt.Println("Using default loglevel: info")
		return
	}
}
