package internal

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigPath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configPath := homePath + "/" + configFileName
	return configPath, nil
}

func Read() (Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return Config{}, err
	}

	cfg := Config{}
	file, err := os.Open(configPath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func SetUser(config Config, userName string) error {
	config.CurrentUserName = userName

	err := writeConfig(config)
	if err != nil {
		return err
	}

	return nil
}

func writeConfig(config Config) error {
	jsonBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, jsonBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
