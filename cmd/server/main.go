package main

import (
	"github.com/btieyris-pixel/notification-service/internal/config"
	"github.com/btieyris-pixel/notification-service/internal/logger"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.ServiceName)

	log.Info("notification-service starting")

	for {
	}
}
