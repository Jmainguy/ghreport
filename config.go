package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config : configuration for what orgs / users / repos and token to use
type Config struct {
	AutoDiscover struct {
		Organizations []struct {
			Name  string `yaml:"name"`
			Topic string `yaml:"topic"`
		} `yaml:"organizations"`
		Users []struct {
			Name  string `yaml:"name"`
			Topic string `yaml:"topic"`
		} `yaml:"users"`
	} `yaml:"autoDiscover"`
	SubscribedRepos []string `yaml:"subscribedRepos"`
	Token           string   `yaml:"token"`
	DefaultOutput   string   `yaml:"defaultOutput"` // new field
}

func ensureConfigDirExists(configDir string) error {
	// Check if config directory exists
	_, err := os.Stat(configDir)
	if os.IsNotExist(err) {
		// Config directory doesn't exist, create it
		err := os.MkdirAll(configDir, 0700) // 0700 means only the owner can read, write, and execute
		if err != nil {
			return fmt.Errorf("failed to create config directory: %v", err)
		}
	} else if err != nil {
		// Some error occurred while checking the existence of the directory
		return fmt.Errorf("failed to check config directory: %v", err)
	}

	return nil
}

func ensureConfigFileExists(configFilePath string) error {
	_, err := os.Stat(configFilePath)
	if os.IsNotExist(err) {
		// Config file doesn't exist, create it with default values
		defaultConfig := Config{
			SubscribedRepos: []string{},
			Token:           "",
			DefaultOutput:   "", // default to empty
		}
		configBytes, err := yaml.Marshal(defaultConfig)
		if err != nil {
			return fmt.Errorf("failed to marshal default config: %v", err)
		}
		err = os.WriteFile(configFilePath, configBytes, 0600) // 0600 means only the owner can read and write
		if err != nil {
			return fmt.Errorf("failed to create config file: %v", err)
		}
	} else if err != nil {
		// Some error occurred while checking the existence of the file
		return fmt.Errorf("failed to check config file: %v", err)
	}

	return nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(homeDir, ".config/ghreport")
	if err := ensureConfigDirExists(configDir); err != nil {
		return "", err
	}
	configFilePath := filepath.Join(configDir, "config.yaml")
	if err := ensureConfigFileExists(configFilePath); err != nil {
		return "", err
	}
	return configFilePath, nil
}

func readConfigFile(configFilePath string) (*Config, error) {
	file, err := os.Open(configFilePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "error closing config file: %v\n", cerr)
		}
	}()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func getConfig() (*Config, error) {
	// Check if config file exists
	configFilePath, err := getConfigFilePath()
	if err == nil {
		if _, err := os.Stat(configFilePath); err == nil {
			// Config file exists, read values from there
			config, err := readConfigFile(configFilePath)
			if err == nil {
				return config, nil
			}
		}
	}

	// If config file doesn't exist or there was an error reading it, fallback to environment variables
	var config Config
	envSubscribedRepos := os.Getenv("subscribedRepos")
	if envSubscribedRepos == "" {
		return &config, fmt.Errorf("env variable subscribedRepos is not defined")
	}
	config.SubscribedRepos = strings.Split(envSubscribedRepos, " ")
	envToken := os.Getenv("ghreportToken")
	if envToken == "" {
		return &config, fmt.Errorf("env variable ghreportToken is not defined")
	}
	config.Token = envToken
	config.DefaultOutput = os.Getenv("ghreportDefaultOutput") // allow override via env

	return &config, nil

}
