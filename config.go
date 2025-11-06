package main

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type UserConfig struct {
	Solved [12]int `json:"solved"`
}

type Prize struct {
	Cents     int
	Solutions int
	Label     string
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

func GetPrizesPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".queens/prizes.txt"
	}
	return filepath.Join(home, ".queens", "prizes.txt")
}

func CreateDefaultPrizes() error {
	prizesPath := GetPrizesPath()

	dir := filepath.Dir(prizesPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	defaultContent := `020,1,Find one solution
020,2,Find two solutions
050,5,Find 5   solutions
100,7,Find 7   solutions
300,12,Find 12  solutions
`

	return os.WriteFile(prizesPath, []byte(defaultContent), 0644)
}

func LoadPrizes() ([]Prize, error) {
	prizesPath := GetPrizesPath()

	file, err := os.Open(prizesPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err := CreateDefaultPrizes(); err != nil {
				return nil, err
			}
			file, err = os.Open(prizesPath)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	defer file.Close()

	var prizes []Prize
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ",", 3)
		if len(parts) != 3 {
			continue
		}

		cents, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			continue
		}

		solutions, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			continue
		}

		label := parts[2]

		prizes = append(prizes, Prize{
			Cents:     cents,
			Solutions: solutions,
			Label:     label,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return prizes, nil
}
