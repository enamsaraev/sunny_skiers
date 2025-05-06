package main

import (
	"os"
	"path/filepath"
	"skiers/internal/config"
	"skiers/internal/events"
	"skiers/pkg/logger"
)

func main() {
	logger.CreateLogger()

	cfg, err := config.LoadConfig(filepath.Join("data", "config.json"))
	if err != nil {
		logger.GetLogger().Errorf("error while loading config: %v", err)
		os.Exit(1)
	}

	eventData, err := events.ParseEventFile(filepath.Join("data", "events"))
	if err != nil {
		logger.GetLogger().Errorf("error while parsing event file: %v", err)
		os.Exit(1)
	}

	err = events.LogCompetitorsData(eventData)
	if err != nil {
		logger.GetLogger().Errorf("error while creating result table: %v", err)
		os.Exit(1)
	}
	events.CreateResultTable(eventData, cfg)
}
