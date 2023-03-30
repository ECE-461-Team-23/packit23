package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds the configuration data
type Config struct {
	LogLevel    string
	LogFile     string
	GithubToken string
}

// NewConfig creates a new Config and reads the environment variables
func NewConfig() *Config {
	viper.SetDefault("LOG_LEVEL", "0")
	viper.SetDefault("GITHUB_TOKEN", "")
	viper.SetDefault("LOG_FILE", "logfile.log")
	viper.AutomaticEnv()

	return &Config{
		LogLevel:    viper.GetString("LOG_LEVEL"),
		LogFile:     viper.GetString("LOG_FILE"),
		GithubToken: viper.GetString("GITHUB_TOKEN"),
	}
}

// Check Github token is set
func (c *Config) CheckToken() error {
	if c.GithubToken == "" {
		return fmt.Errorf("GITHUB_TOKEN not set")
	}

	return nil
}

// Check if log variables are set
func (c *Config) CheckLog() string {
	return fmt.Sprintf("Current LOG_LEVEL: %s, Current LOG_FILE output location: %s", c.LogLevel, c.LogFile)
}
