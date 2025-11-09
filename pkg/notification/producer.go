package notification

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EmailPayload struct {
	Recipient string `json:"recipient"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}
type Payload struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
	Type    string `json:"type"`
	Email   string `json:"email"`
}

type Notification struct {
	Type    string          `json:"type"`    // "email", "sms", "push"
	Payload json.RawMessage `json:"payload"` // EmailPayload, SMSPayload и т.д.
}

type Producer interface {
	Send(ctx context.Context, notification Notification) error
	Close() error
}

type RabbitMQProducer struct {
	channel   *amqp.Channel
	queueName string
}

func NewRabbitMQProducer(amqpURL, queueName string) (*RabbitMQProducer, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return &RabbitMQProducer{
		channel:   ch,
		queueName: queueName,
	}, nil
}

func (p *RabbitMQProducer) Send(ctx context.Context, notification Notification) error {
	body, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	return p.channel.PublishWithContext(ctx,
		"",          // exchange
		p.queueName, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		})
}

func (p *RabbitMQProducer) Close() error {
	return p.channel.Close()
}
func NewEmailNotification(to, subject, body string) (Notification, error) {
	payload, err := json.Marshal(EmailPayload{
		Recipient: to,
		Subject:   subject,
		Body:      body,
	})
	if err != nil {
		return Notification{}, fmt.Errorf("failed to marshal email payload: %w", err)
	}

	return Notification{
		Type:    "email",
		Payload: payload,
	}, nil
}
