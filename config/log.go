package config

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Global variables for logging
var logFile *os.File
var logger *zerolog.Logger
var traceID string

// InitiateLog sets up the logger, writes to the log file, and adds a trace ID to the context
func InitiateLog() error {
	traceID = uuid.New().String()

	if _, err := os.Stat("log"); os.IsNotExist(err) {
		err := os.Mkdir("log", os.ModePerm)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create log directory")
			return fmt.Errorf("failed to create log directory: %v", err)
		}
	}

	var err error
	logFile, err = os.OpenFile("log/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open log file")
		return fmt.Errorf("failed to open log file: %v", err)
	}

	// Set up zerolog with trace ID and timestamp
	createdLog := zerolog.New(logFile).With().
		Str("trace_id", traceID).
		Timestamp().
		Logger()

	// Store logger as a pointer
	logger = &createdLog

	logger.Info().Msg("Zerolog initiated successfully")

	return nil
}

func GetLog() *zerolog.Logger {
	return logger
}

func SetTraceID(newTraceID string) {
	traceID = newTraceID

	// Rebuild the logger with the updated trace ID
	createdLog := zerolog.New(logFile).With().
		Str("trace_id", traceID).
		Timestamp().
		Logger()

	logger = &createdLog
}

func GetTraceID() string {
	return traceID
}

// CloseLog closes the log file when the application shuts down
func CloseLog() {
	if logFile != nil {
		err := logFile.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close log file")
		}
	}
}
