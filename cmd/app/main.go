package main

import (
	"fmt"
	"path/filepath"
	"skiers/internal/events"
)

func main() {
	/*
		cfg, err := config.LoadConfig(filepath.Join("data", "config.json"))
		if err != nil {
			fmt.Println(err)
		}
	*/

	err := events.LogCompetitorsData(filepath.Join("data", "events"))
	if err != nil {
		fmt.Println(err)
	}
}
