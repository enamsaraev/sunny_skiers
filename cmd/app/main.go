package main

import (
	"fmt"
	"path/filepath"
	"skiers/internal/config"
	"skiers/internal/events"
)

func main() {
	cfg, err := config.LoadConfig(filepath.Join("data", "config.json"))
	if err != nil {
		fmt.Println(err)
	}

	eventData := events.ParseEventFile(filepath.Join("data", "events"))
	err = events.LogCompetitorsData(eventData)
	if err != nil {
		fmt.Println(err)
	}
	events.CreateResultTable(eventData, cfg)
}
