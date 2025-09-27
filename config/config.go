package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Payment  PaymentConfig
	Storage  StorageConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Env          string
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	SSLMode  string
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	SecretKey     string
	Expiration    time.Duration
	RefreshSecret string
}

// PaymentConfig holds payment-related configuration
type PaymentConfig struct {
	StripeSecretKey      string
	StripePublishableKey string
	PayPalClientID       string
	PayPalSecret         string
}

// StorageConfig holds file storage configuration
type StorageConfig struct {
	CloudinaryURL string
	UploadPath    string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	readTimeout := getEnvAsDuration("SERVER_READ_TIMEOUT", 10*time.Second)
	writeTimeout := getEnvAsDuration("SERVER_WRITE_TIMEOUT", 10*time.Second)

	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			Env:          getEnv("GO_ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			Username: getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Database: getEnv("DB_NAME", "health-store-db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			SecretKey:     getEnv("JWT_SECRET_KEY", "your-super-secret-jwt-key-change-this-in-production-2024"),
			Expiration:    getEnvAsDuration("JWT_EXPIRATION", 24*time.Hour),
			RefreshSecret: getEnv("JWT_REFRESH_SECRET", "your-refresh-secret-key"),
		},
		Payment: PaymentConfig{
			StripeSecretKey:      getEnv("STRIPE_SECRET_KEY", ""),
			StripePublishableKey: getEnv("STRIPE_PUBLISHABLE_KEY", ""),
			PayPalClientID:       getEnv("PAYPAL_CLIENT_ID", ""),
			PayPalSecret:         getEnv("PAYPAL_SECRET", ""),
		},
		Storage: StorageConfig{
			CloudinaryURL: getEnv("CLOUDINARY_URL", ""),
			UploadPath:    getEnv("UPLOAD_PATH", "./uploads"),
		},
	}
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsInt gets an environment variable as integer with a fallback value
func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Printf("Warning: Invalid integer value for %s: %s, using fallback: %d", key, value, fallback)
	}
	return fallback
}

// getEnvAsDuration gets an environment variable as duration with a fallback value
func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
		log.Printf("Warning: Invalid duration value for %s: %s, using fallback: %v", key, value, fallback)
	}
	return fallback
}

// getEnvAsBool gets an environment variable as boolean with a fallback value
func getEnvAsBool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
		log.Printf("Warning: Invalid boolean value for %s: %s, using fallback: %t", key, value, fallback)
	}
	return fallback
}

// IsDevelopment returns true if the environment is development
func (c *Config) IsDevelopment() bool {
	return c.Server.Env == "development"
}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	return c.Server.Env == "production"
}

// IsTest returns true if the environment is test
func (c *Config) IsTest() bool {
	return c.Server.Env == "test"
}

// GetDatabaseDSN returns the database connection string
func (c *Config) GetDatabaseDSN() string {
	return c.Database.Username + ":" + c.Database.Password + "@tcp(" + c.Database.Host + ":" + c.Database.Port + ")/" + c.Database.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
}
