package rabbit

import (
	"context"
	"encoding/json"

	"go-boilerplate/config"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

type handler func([]byte)

type pusher struct {
	Channel  *amqp.Channel
	handlers map[string]handler
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
	pusher.handlers = map[string]handler{}
	Service = pusher

	log.Info().Msgf("[AMQP] server is running")
}

func (p *pusher) Publish(ctx context.Context, queueName string, data interface{}) {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Panic().Err(err)
	}
	err = p.Channel.PublishWithContext(
		context.Background(),
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(bytes),
		},
	)

	failOnError(err, "Failed to publish a message")

	log.Info().Msgf("[AMQP] Sent %s\n", data)
}

func (p *pusher) RegisterConsumer(queueName string, callback func([]byte)) {
	log.Info().Msgf("[AMQP] Register consumer: %s", queueName)
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
		log.Panic().Err(err)
	}
	go func() {
		for message := range messages {
			log.Printf("[AMQP] queue %s - received message: %s\n", queueName, message.Body)
			callback(message.Body)
		}
	}()
}

func (p *pusher) RegisterConsumers(consumers []string) {
	for _, consumer := range consumers {
		var handler handler
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
		log.Panic().Err(err).Msg(msg)
	}
}
