package rabbit

import (
	"context"
	"encoding/json"
	"fmt"

	"go-boilerplate/config"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

type handlerFunc func([]byte)

type pusher struct {
	Channel  *amqp.Channel
	handlers map[string]handlerFunc
}

var Service *pusher

func Setup(cfg config.RabbitConfig) {
	conn, err := amqp.Dial(cfg.RabbitHost)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()

	failOnError(err, "Failed to open a channel")

	pusher := new(pusher)
	pusher.Channel = ch
	// Feel free to add handlers
	pusher.handlers = map[string]handlerFunc{}
	Service = pusher

	slog.Info("[AMQP]", "message", "server is running")
}

func (p *pusher) Publish(_ context.Context, queueName string, data any) {
	bytes, err := json.Marshal(data)
	if err != nil {
		slog.Error("[AMQP]", "Publish failed", slog.String("error", err.Error()))
		panic(err)
	}

	err = p.Channel.PublishWithContext(
		context.Background(),
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bytes,
		},
	)

	failOnError(err, "Failed to publish a message")

	slog.Info("[AMQP]", "message", fmt.Sprintf("Sent %s", data))
}

func (p *pusher) RegisterConsumer(queueName string, callback func([]byte)) {
	slog.Info("[AMQP]", "message", fmt.Sprintf("Register consumer: %s", queueName))
	messages, err := p.Channel.Consume(
		queueName, // queue name
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no local
		false,     // no wait
		nil,       // arguments
	)
	if err != nil {
		slog.Error("[AMQP]", "RegisterConsumer failed", slog.String("error", err.Error()))
		panic(err)
	}
	go func() {
		for message := range messages {
			slog.Info("[AMQP]", "message", fmt.Sprintf("queue %s - received message: %s", queueName, message.Body))
			callback(message.Body)
		}
	}()
}

func (p *pusher) RegisterConsumers(consumers []string) {
	for _, consumer := range consumers {
		var handler handlerFunc
		if handler = p.handlers[consumer]; handler != nil {
			p.RegisterConsumer(consumer, handler)
		}
	}
}

func (p *pusher) CloseChannel() {
	p.Channel.Close()
}

func failOnError(err error, msg string) {
	if err != nil {
		slog.Error("[AMQP]", msg, slog.String("error", err.Error()))
		panic(err)
	}
}
