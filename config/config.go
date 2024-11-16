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

	configInstance = config

	return config, nil
}

func Get() *Config {
	if configInstance == nil {
		_, err := LoadConfig("config.json")
		if err != nil {
			log.Fatal("failed to load config file	", zap.Error(err))
		}
	}
	return configInstance
}
