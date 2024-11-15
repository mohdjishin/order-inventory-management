package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	log "github.com/mohdjishin/order-inventory-management/logger"
	"github.com/rs/zerolog"
)

type Config struct {
	DSN      string `json:"dsn"`
	LogLevel string `json:"LogLevel"`
	LogFile  string `json:"LogFile"`
	Port     string `json:"port"`
}

var configInstance *Config

// LoadConfig loads configuration from a config.json file
func LoadConfig(filePath string) (*Config, error) {
	if configInstance != nil {
		return configInstance, nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, fmt.Errorf("could not decode config file: %w", err)
	}

	// Set the config instance for future access
	configInstance = config

	// Initialize logger based on config
	initLogger(config.LogLevel, config.LogFile)

	return config, nil
}

// GetConfig returns the loaded config instance
func Get() *Config {
	if configInstance == nil {
		_, err := LoadConfig("config.json") // Default path to the config file
		if err != nil {
			log.Fatal().Err(err).Msg("failed to load config file	")
		}
	}
	return configInstance
}

// initLogger initializes the logger based on log level and log file
func initLogger(logLevel, logFile string) {
	var level zerolog.Level
	switch strings.ToLower(logLevel) {
	case "debug":
		level = zerolog.DebugLevel
	case "info":
		level = zerolog.InfoLevel
	case "warn":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	case "fatal":
		level = zerolog.FatalLevel
	default:
		level = zerolog.InfoLevel
	}

	// Set up the logger to log to a file or console
	if logFile != "" {
		// Log to file
		file, err := os.Create(logFile)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create log file")
		}
		zerolog.SetGlobalLevel(level)
		log.Logger = zerolog.New(file).With().Timestamp().Logger()
	} else {
		// Log to console
		zerolog.SetGlobalLevel(level)
		log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	}
}
