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
	if err := WriteConfig(c); err != nil {
		return err
	}
	return nil

}

func WriteConfig(config *Config) error {
	configFile, err := getConfigFilePath()
	if err != nil {
		return err
	}
	jsonData, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("cannot convert to json")
	}
	file, err := os.Create(configFile)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	// Write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}
	return nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	configFile := homeDir + "/.gatorconfig.json"
	if err != nil {
		return "", err
	}
	return configFile, nil
}
