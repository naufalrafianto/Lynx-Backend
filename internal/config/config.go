package config

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
	}
	Redis struct {
		Host     string
		Port     string
		Password string
	}
	JWT struct {
		Secret          string
		ExpirationHours int
	}
	SMTP struct {
		Host     string
		Port     string
		Email    string
		Password string
	}
	RateLimit struct {
		Requests int
		Duration time.Duration
	}
	Password struct {
		MinLength        int
		RequireSpecial   bool
		RequireNumbers   bool
		RequireUppercase bool
	}
	OTP struct {
		Length            int
		ExpirationMinutes int
		MaxAttempts       int
	}
	App struct {
		Environment string
		Port        string
		GinMode     string
	}
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, errors.New("error loading .env file")
	}

	config := &Config{}

	// Database config
	config.Database.Host = getEnv("DB_HOST", "")
	config.Database.Port = getEnv("DB_PORT", "5432")
	config.Database.User = getEnv("DB_USER", "")
	config.Database.Password = getEnv("DB_PASSWORD", "")
	config.Database.DBName = getEnv("DB_NAME", "")

	// Redis config
	config.Redis.Host = getEnv("REDIS_HOST", "")
	config.Redis.Port = getEnv("REDIS_PORT", "6379")
	config.Redis.Password = getEnv("REDIS_PASSWORD", "")

	// JWT config
	config.JWT.Secret = getEnv("JWT_SECRET", "")
	config.JWT.ExpirationHours = getEnvAsInt("JWT_EXPIRATION_HOURS", 24)

	// SMTP config
	config.SMTP.Host = getEnv("SMTP_HOST", "")
	config.SMTP.Port = getEnv("SMTP_PORT", "")
	config.SMTP.Email = getEnv("SMTP_EMAIL", "")
	config.SMTP.Password = getEnv("SMTP_PASSWORD", "")

	// Rate limiting config
	config.RateLimit.Requests = getEnvAsInt("RATE_LIMIT_REQUESTS", 100)
	config.RateLimit.Duration = time.Duration(getEnvAsInt("RATE_LIMIT_DURATION", 60)) * time.Second

	// Password requirements
	config.Password.MinLength = getEnvAsInt("MIN_PASSWORD_LENGTH", 12)
	config.Password.RequireSpecial = getEnvAsBool("REQUIRE_SPECIAL_CHARS", true)
	config.Password.RequireNumbers = getEnvAsBool("REQUIRE_NUMBERS", true)
	config.Password.RequireUppercase = getEnvAsBool("REQUIRE_UPPERCASE", true)

	// OTP settings
	config.OTP.Length = getEnvAsInt("OTP_LENGTH", 6)
	config.OTP.ExpirationMinutes = getEnvAsInt("OTP_EXPIRATION_MINUTES", 5)
	config.OTP.MaxAttempts = getEnvAsInt("MAX_OTP_ATTEMPTS", 3)

	// Application settings
	config.App.Environment = getEnv("ENV", "development")
	config.App.Port = getEnv("PORT", "8080")
	config.App.GinMode = getEnv("GIN_MODE", "release")

	return config, validateConfig(config)
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := strconv.ParseBool(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

func validateConfig(config *Config) error {
	if config.JWT.Secret == "" {
		return errors.New("JWT_SECRET is required")
	}
	if config.Database.Password == "" {
		return errors.New("DB_PASSWORD is required")
	}
	// Add more validation as needed
	return nil
}
