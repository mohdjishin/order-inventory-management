package config

import (
	"encoding/json"
	"fmt"
	"os"

	log "github.com/mohdjishin/order-inventory-management/logger"
	"go.uber.org/zap"
)

type Config struct {
	DSN      string `json:"dsn"`
	LogLevel string `json:"LogLevel"`
	LogFile  string `json:"LogFile"`
	Port     string `json:"port"`
	JwtKey   string `json:"jwtSecret"`
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
	// initLogger(config.LogLevel, config.LogFile)

	return config, nil
}

// GetConfig returns the loaded config instance
func Get() *Config {
	if configInstance == nil {
		_, err := LoadConfig("config.json") // Default path to the config file
		if err != nil {
			log.Fatal("failed to load config file	", zap.Error(err))
		}
	}
	return configInstance
}

// initLogger initializes the logger based on log level and log file
