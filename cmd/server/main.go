package main

import (
	"context"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	awssqs "github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/btieyris-pixel/notification-service/internal/config"
	"github.com/btieyris-pixel/notification-service/internal/db"
	"github.com/btieyris-pixel/notification-service/internal/fcm"
	"github.com/btieyris-pixel/notification-service/internal/logger"
	"github.com/btieyris-pixel/notification-service/internal/repository"
	internalsqs "github.com/btieyris-pixel/notification-service/internal/sqs"
)

// Handler de eventos SQS:
// - Obtiene token del driver desde PostgreSQL
// - Envía push vía FCM
// - Si falla, NO elimina mensaje (retry por SQS)
func main() {
	ctx := context.Background()

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

	log.Info("fcm client initialized")

	awsCfg, err := awsconfig.LoadDefaultConfig(
		ctx,
		awsconfig.WithRegion(cfg.AWSRegion),
	)
	if err != nil {
		log.Error(err.Error())
		log.Fatal("failed to load aws config")
	}

	sqsClient := awssqs.NewFromConfig(awsCfg)

	consumer := internalsqs.New(
		sqsClient,
		cfg.SQSQueueURL,
		int32(cfg.SQSWaitTimeSeconds),
		int32(cfg.SQSMaxMessages),
		func(ctx context.Context, event internalsqs.Event) error {
			token, err := repo.GetDriverToken(ctx, event.DriverID)
			if err != nil {
				log.Error("token not found for driver: " + event.DriverID)
				return err
			}

			_, err = fcmClient.SendToToken(ctx, fcm.PushMessage{
				Token:   token,
				Title:   "Nueva orden",
				Body:    "Tienes una nueva orden disponible",
				OrderID: event.OrderID,
				Event:   event.EventType,
			})
			if err != nil {
				log.Error("push failed after retries: " + err.Error())
				return err
			}

			log.Info("push sent for order: " + event.OrderID)
			return nil
		},
	)

	log.Info("sqs consumer starting")

	if err := consumer.Start(ctx); err != nil {
		log.Error(err.Error())
		log.Fatal("sqs consumer stopped")
	}
}
