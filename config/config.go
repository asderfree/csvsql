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

func (c *Config) isInMemoryDB() bool {
	return c.DatabasePath == ":memory:"
}

var Gcfg *Config

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

// if the Gcfg's path is not exist, just create this file:
func creatDbFile(f string) {
	info, err := os.Stat(f)
	if err == nil && !info.IsDir() {
		// File exists and is not a directory, nothing to do
		return
	}
	if os.IsNotExist(err) {
		file, err := os.Create(f)
		if err != nil {
			// Could not create file, optionally handle error
			return
		}
		file.Close()
	}
	// If err is not nil and not IsNotExist, do nothing (could log if verbose)
	// If file exists but is a directory, do nothing
	// FIX if the file already exists, try to exit the program and panic the message.
}

func init() {
	Gcfg = Load()
	if !Gcfg.isInMemoryDB() {
		creatDbFile(Gcfg.DatabasePath)
	}
}
