package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Laps        int    `json:"laps"`
	LapLen      int    `json:"lapLen"`
	PenaltyLen  int    `json:"penaltyLen"`
	FiringLines int    `json:"firingLines"`
	Start       string `json:"start"`
	StartDelta  string `json:"startDelta"`
}

func LoadConfig(path string) (*Config, error) {
	// реализация чтения и парсинга JSON
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error while reading a file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(file, &cfg); err != nil {
		return nil, fmt.Errorf("error while parsing json: %w", err)
	}

	// validation
	if cfg.Laps <= 0 {
		return nil, fmt.Errorf("number of laps must be positive")
	}

	return &cfg, nil
}
