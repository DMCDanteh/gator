package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}
	location := filepath.Join(home, ".gatorconfig.json")
	data, err := os.ReadFile(location)
	if err != nil {
		return Config{}, err
	}
	cfg := Config{}
	if err = json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName
	return write(*cfg)
}

func write(cfg Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	location := filepath.Join(home, ".gatorconfig.json")
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(location, data, 0600)
}
