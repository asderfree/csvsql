package config

import (
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	DatabasePath string
	MaxFileSize  int64
	Verbose      bool
}

// Load loads configuration from environment variables and defaults
func Load() *Config {
	config := &Config{
		DatabasePath: ":memory:",        // Default to in-memory database
		MaxFileSize:  100 * 1024 * 1024, // 100MB default
		Verbose:      false,
	}

	// Override with environment variables if set
	if dbPath := os.Getenv("DANA_DB_PATH"); dbPath != "" {
		config.DatabasePath = dbPath
	}

	if maxSize := os.Getenv("DANA_MAX_FILE_SIZE"); maxSize != "" {
		if size, err := strconv.ParseInt(maxSize, 10, 64); err == nil {
			config.MaxFileSize = size
		}
	}

	if verbose := os.Getenv("DANA_VERBOSE"); verbose == "true" {
		config.Verbose = true
	}

	return config
}
