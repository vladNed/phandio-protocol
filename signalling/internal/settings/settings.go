package settings

import (
	"os"

	"github.com/joho/godotenv"
)

type RedisSettings struct {
	Address  string
	Password string
	DB       int
}

func NewRedisSettings() *RedisSettings {
	redisSettings := &RedisSettings{
		Address:  os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	}

	// Check if address and password are set
	if redisSettings.Address == "" {
		redisSettings.Address = "localhost:6379"
	}
	if redisSettings.Password == "" {
		redisSettings.Password = ""
	}

	return redisSettings
}

type DefaultSettings struct {
	// Port to listen on
	Port string

	// Host to listen on
	Host string

	// Log level
	LogLevel string

	// Redis settings
	Redis *RedisSettings

	// Certs
	CertFile string
	KeyFile  string
}

func GetSettings() DefaultSettings {
	godotenv.Load()

	redisSettings := NewRedisSettings()
	defaultSettings := DefaultSettings{
		Port:     "8080",
		Host:     "0.0.0.0",
		LogLevel: "info",
		Redis:    redisSettings,
		CertFile: "",
		KeyFile:  "",
	}
	if port := os.Getenv("PORT"); port != "" {
		defaultSettings.Port = port
	}
	if host := os.Getenv("HOST"); host != "" {
		defaultSettings.Host = host
	}
	if certFile := os.Getenv("CERT_FILE"); certFile != "" {
		defaultSettings.CertFile = certFile
	}
	if keyFile := os.Getenv("KEY_FILE"); keyFile != "" {
		defaultSettings.KeyFile = keyFile
	}

	return defaultSettings
}

func (s *DefaultSettings) GetAddress() string {
	return s.Host + ":" + s.Port
}
