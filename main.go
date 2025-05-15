package main

import (
	"log"
	"submit_meter_readings/bot"
	"submit_meter_readings/config"
	"submit_meter_readings/internal/storage"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Config error:", err)
	}

	pgStorage, err := storage.NewPostgresStorage(cfg)
	if err != nil {
		log.Fatal("Storage init error:", err)
	}
	defer pgStorage.Close()

	bot, err := bot.NewBot(cfg, pgStorage)
	if err != nil {
		log.Fatal("Bot init error:", err)
	}

	bot.Start()

	select {}
}
