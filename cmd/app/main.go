package main

import (
	"fmt"
	"skiers/internal/config"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		panic(err)
	}

	fmt.Println(cfg)
}
