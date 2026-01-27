package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Environment string           `json:"environment"`
	Address     string           `json:"address"`
	Encryption  EncryptionConfig `json:"encryption"`
	Database    DatabaseConfig   `json:"database"`
	Polar       PolarConfig      `json:"polar"`
}

type PolarConfig struct {
	WebhookSecret string `json:"webhookSecret"`
}

type EncryptionConfig struct {
	Key string `json:"key"`
}

type DatabaseConfig struct {
	ConnectionString string `json:"connectionString"`
}

// Load reads configuration from environment variables and optionally from a config file
func Load() (*Config, error) {
	// Start with defaults
	config := setDefaults()

	// Override with file config if present
	if configPath := os.Getenv("CONFIG_FILE"); configPath != "" {
		if err := loadFromFile(configPath, config); err != nil {
			return nil, fmt.Errorf("error loading config file: %w", err)
		}
	}

	// Override with environment variables
	loadFromEnv(config)

	// Validate final configuration
	if err := validate(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

func loadFromFile(path string, config *Config) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(config)
}

func loadFromEnv(config *Config) {
	if environment := os.Getenv("ENVIRONMENT"); environment != "" {
		config.Environment = environment
	}

	if address := os.Getenv("ADDRESS"); address != "" {
		config.Address = address
	}

	// Encryption configuration
	if encryptionKey := os.Getenv("ENCRYPTION_KEY"); encryptionKey != "" {
		config.Encryption.Key = encryptionKey
	}

	// Database configuration
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		config.Database = DatabaseConfig{
			ConnectionString: databaseURL,
		}
	}

	// Polar configuration
	if polarWebhookSecret := os.Getenv("POLAR_WEBHOOK_SECRET"); polarWebhookSecret != "" {
		config.Polar.WebhookSecret = polarWebhookSecret
	}
}

func setDefaults() *Config {
	return &Config{
		Environment: "dev",
		Address:     ":8080",
	}
}

func validate(config *Config) error {
	if config.Environment == "" {
		return fmt.Errorf("environment is required")
	}

	if config.Address == "" {
		return fmt.Errorf("address is required")
	}

	// Encryption key validation
	if config.Encryption.Key != "" {
		// Decode base64 key and check decoded length
		decodedKey, err := base64.StdEncoding.DecodeString(config.Encryption.Key)
		if err != nil {
			return fmt.Errorf("encryption key must be valid base64: %w", err)
		}
		if len(decodedKey) != 32 {
			return fmt.Errorf("encryption key must decode to exactly 32 bytes (256 bits), got %d bytes", len(decodedKey))
		}
	}

	// Database configuration validation
	if config.Database.ConnectionString == "" {
		return fmt.Errorf("database connection string is required")
	}

	return nil
}
