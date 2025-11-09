package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"

	"github.com/prankevich/Auth_service/internal/bootstrap"
	"github.com/prankevich/Auth_service/internal/config"
	"github.com/prankevich/Auth_service/pkg/notification"
)

func main() {

	_ = godotenv.Load(".env")

	var cfg config.Config
	err := envconfig.ProcessWith(context.TODO(), &envconfig.Config{
		Target:   &cfg,
		Lookuper: envconfig.OsLookuper(),
	})
	if err != nil {
		log.Fatalf(" Ошибка загрузки конфигурации: %v", err)
	}
	amqpURL := "amqp://guest:guest@localhost:5672/"
	queueName := "notifications"

	producer, err := notification.NewRabbitMQProducer(amqpURL, queueName)
	if err != nil {
		log.Fatalf(" Ошибка инициализации RabbitMQ producer: %v", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Printf(" ️ Ошибка при закрытии RabbitMQ producer: %v", err)
		}
	}()

	// Обработка сигнала завершения
	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, os.Interrupt)

	app := bootstrap.New(cfg, producer)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-quitSignal
		cancel()
	}()

	app.Run(ctx)
}
