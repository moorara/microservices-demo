package config

import (
	"io/ioutil"
	"os"
)

const (
	defaultLogLevel    = "info"
	defaultServiceName = "go-service"
	defaultServicePort = ":4010"
	defaultRedisURL    = "redis://localhost:6379"
)

// Config represents configurations of service
type Config struct {
	LogLevel    string
	ServiceName string
	ServicePort string
	RedisURL    string
}

func getValue(name, defaultValue string) string {
	// Try reading from environment variable directly
	value := os.Getenv(name)
	if value != "" {
		return value
	}

	// Try reading from a file specified by environment variable
	filepath := os.Getenv(name + "_FILE")
	if filepath != "" {
		data, err := ioutil.ReadFile(filepath)
		if err != nil {
			panic(err)
		}
		return string(data)
	}

	return defaultValue
}

// GetConfig retrieves configuratinos
func GetConfig() Config {
	return Config{
		LogLevel:    getValue("LOG_LEVEL", defaultLogLevel),
		ServiceName: getValue("SERVICE_NAME", defaultServiceName),
		ServicePort: getValue("SERVICE_PORT", defaultServicePort),
		RedisURL:    getValue("REDIS_URL", defaultRedisURL),
	}
}
