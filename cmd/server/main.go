package main

import (
	"context"

	"github.com/btieyris-pixel/notification-service/internal/config"
	"github.com/btieyris-pixel/notification-service/internal/db"
	"github.com/btieyris-pixel/notification-service/internal/fcm"
	"github.com/btieyris-pixel/notification-service/internal/logger"
	"github.com/btieyris-pixel/notification-service/internal/repository"
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
	repo := repository.NewTokenRepository(pool)

	token, err := repo.GetDriverToken(context.Background(), "driver-123")
	if err != nil {
		log.Error("failed to get driver token: " + err.Error())
	} else {
		log.Info("token found: " + token)
	}

	ctx := context.Background()

	fcmClient, err := fcm.New(
		ctx,
		cfg.FCMEnabled,
		cfg.FCMProjectID,
		cfg.FCMCredentialsJSON,
	)

	if err != nil {
		log.Error(err.Error())
		log.Fatal("failed to init fcm client")
	}

	_, err = fcmClient.SendToToken(ctx, fcm.PushMessage{
		Token:   token,
		Title:   "Nueva orden",
		Body:    "Tienes una nueva orden disponible",
		OrderID: "test-order-123",
		Event:   "NEW_ORDER",
	})

	if err != nil {
		log.Error("failed to send push: " + err.Error())
	} else {
		log.Info("push sent successfully")
	}
	for {
	}
}
