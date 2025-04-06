package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Db_url       string `json:"db_url"`
	Current_user string `json:"current_user_name"`
}

func ReadConfig() (Config, error) {
	configFile, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	config_content, err := os.ReadFile(configFile)
	if err != nil {
		return Config{}, err
	}
	var config Config
	if err := json.Unmarshal(config_content, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func (c *Config) SetUser() error {
	username := os.Getenv("USER") // For Unix/Linux systems
	if username == "" {
		username = os.Getenv("USERNAME") // For Windows systems
	}
	if username == "" {
		return fmt.Errorf("no username found")
	}
	c.Current_user = username
	return nil

}

func writeConfig(config Config) error {

}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	configFile := homeDir + "/.gatorconfig.json"
	if err != nil {
		return "", err
	}
	return configFile, nil
}
