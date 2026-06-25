package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"voice-bridge-client/internal/client"
)

type Config struct {
	ServerURL string    `json:"serverUrl"`
	Token     string    `json:"token"`
	UserID    string    `json:"userId"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func Load() (Config, error) {
	path, err := configPath()
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return Config{}, nil
	}
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	return cfg, json.Unmarshal(data, &cfg)
}

func Save(cfg Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func Clear() error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if os.IsNotExist(os.Remove(path)) {
		return nil
	}
	return nil
}

func (c Config) Valid() bool {
	return c.Token != "" && time.Now().Before(c.ExpiresAt)
}

func (c *Config) SetSession(session client.Session) {
	c.Token = session.Token
	c.UserID = session.UserID
	c.ExpiresAt = session.ExpiresAt
}

func configPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "voice-bridge-client", "config.json"), nil
}
