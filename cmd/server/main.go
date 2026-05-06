package main

import (
	"github.com/btieyris-pixel/notification-service/internal/config"
	"github.com/btieyris-pixel/notification-service/internal/db"
	"github.com/btieyris-pixel/notification-service/internal/logger"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.ServiceName)

	log.Info("notification-service starting")

	pool, err := db.New(cfg.PostgresDSN)
	if err != nil {
		log.Error(err.Error())
		log.Fatal("failed to connect to database")
	}

	defer pool.Close()

	log.Info("database connected")

	for {
	}
}
