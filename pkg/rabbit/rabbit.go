package rabbit

import (
	"fmt"

	"go-boilerplate/config"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

type handlerFunc func([]byte)

type pusher struct {
	chans    []*amqp.Channel
	conn     *amqp.Connection
	handlers map[string]handlerFunc
}

var Service *pusher

func Setup(cfg config.RabbitConfig) {
	conn, err := amqp.Dial(cfg.RabbitHost)
	failOnError(err, "Failed to connect to RabbitMQ")

	failOnError(err, "Failed to open a channel")

	pusher := new(pusher)
	pusher.conn = conn

	Service = pusher

	slog.Info("[AMQP]", "message", "server is running")
}

func (p *pusher) RegisterConsumer(queueName string, callback func([]byte)) {
	slog.Info("[AMQP]", "message", fmt.Sprintf("Register consumer: %s", queueName))

	ch, err := p.conn.Channel()
	if err != nil {
		slog.Error("[AMQP]", "RegisterConsumer failed", slog.String("error", err.Error()))
		panic(err)
	}

	p.chans = append(p.chans, ch)

	messages, err := ch.Consume(
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

func (p *pusher) CloseChannels() error {
	for _, ch := range p.chans {
		err := ch.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		slog.Error("[AMQP]", msg, slog.String("error", err.Error()))
		panic(err)
	}
}
