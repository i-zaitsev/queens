package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type UserConfig struct {
	Solved [12]int `json:"solved"`
}

func GetConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".queens/results.json"
	}
	return filepath.Join(home, ".queens", "results.json")
}

func LoadConfig() (*UserConfig, error) {
	configPath := GetConfigPath()

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &UserConfig{}, nil
		}
		return nil, err
	}

	var config UserConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func SaveConfig(config *UserConfig) error {
	configPath := GetConfigPath()

	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

func (c *UserConfig) MarkSolved(solutionNum int) {
	if solutionNum >= 1 && solutionNum <= 12 {
		c.Solved[solutionNum-1] = 1
	}
}
