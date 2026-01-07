package config

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	SecureCRT string `json:"securecrt"`
	FileZilla string `json:"filezilla"`
}

const cfgFile = "agent.json"

func load() (*Config, error) {
	var c Config
	data, err := os.ReadFile(cfgFile)
	if err != nil {
		return &c, nil
	}
	_ = json.Unmarshal(data, &c)
	return &c, nil
}

func save(c *Config) error {
	data, _ := json.MarshalIndent(c, "", "  ")
	return os.WriteFile(cfgFile, data, 0644)
}

func GetFileZillaPath() (string, error) {
	c, _ := load()
	if c.FileZilla != "" {
		return c.FileZilla, nil
	}
	return "", errors.New("FileZilla not configured")
}
