package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	APIURL   string `json:"api_url"`
	Token    string `json:"token"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

func configDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".uxvalidator")
}

func configPath() string {
	return filepath.Join(configDir(), "config.json")
}

func Load() (*Config, error) {
	data, err := os.ReadFile(configPath())
	if err != nil {
		return &Config{APIURL: "https://api.deltix.ai"}, nil
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return &Config{APIURL: "https://api.deltix.ai"}, nil
	}
	if cfg.APIURL == "" {
		cfg.APIURL = "https://api.deltix.ai"
	}
	return &cfg, nil
}

func Save(cfg *Config) error {
	if err := os.MkdirAll(configDir(), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath(), data, 0600)
}

func (c *Config) IsLoggedIn() bool {
	return c.Token != ""
}
