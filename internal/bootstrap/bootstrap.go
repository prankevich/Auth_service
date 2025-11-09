package bootstrap

import (
	"context"
	"fmt"
	"github.com/prankevich/Auth_service/internal/adapter/driven/dbstore"
	"github.com/prankevich/Auth_service/internal/adapter/driving/http"
	"github.com/prankevich/Auth_service/internal/config"
	"github.com/prankevich/Auth_service/internal/usecase"
	"github.com/prankevich/Auth_service/pkg/notification"
)

func initLayers(cfg config.Config) *App {
	teardown := make([]func(), 0)

	// Инициализация БД
	db, err := initDB(*cfg.Postgres)
	if err != nil {
		panic(err)
	}
	storage := dbstore.New(db)
	teardown = append(teardown, func() {
		if err := db.Close(); err != nil {
			fmt.Println("DB close error:", err)
		}
	})

	// Инициализация RabbitMQ producer
	amqpURL := "amqp://guest:guest@localhost:5672/"
	queueName := "notifications"

	producer, err := notification.NewRabbitMQProducer(amqpURL, queueName)
	if err != nil {
		panic(fmt.Errorf("RabbitMQ producer init error: %w", err))
	}
	teardown = append(teardown, func() {
		if err := producer.Close(); err != nil {
			fmt.Println("RabbitMQ producer close error:", err)
		}
	})

	// Инициализация usecase-слоя
	uc := usecase.New(cfg, storage)

	// Инициализация HTTP-сервера
	httpSrv := http.New(&cfg, uc, producer)
	teardown = append(teardown, func() {
		ctxShutDown, cancel := context.WithTimeout(context.Background(), gracefulDeadline)
		defer cancel()
		if err := httpSrv.Shutdown(ctxShutDown); err != nil {
			fmt.Println("HTTP shutdown error:", err)
		}
	})

	return &App{
		cfg:      cfg,
		rest:     httpSrv,
		useCases: uc,
		teardown: teardown,
	}
}
