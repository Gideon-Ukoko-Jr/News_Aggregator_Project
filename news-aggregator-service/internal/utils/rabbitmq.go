package utils

import (
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"log"
)

// RabbitMQConfig holds the RabbitMQ connection details.
type RabbitMQConfig struct {
	URL      string
	Exchange string
}

type RabbitMQPublisher struct {
	channel  *amqp.Channel
	exchange string
}

// RabbitMQConsumer holds the RabbitMQ consumer details.
type RabbitMQConsumer struct {
	channel  *amqp.Channel
	queue    *amqp.Queue
	messages <-chan amqp.Delivery
}

// NewRabbitMQConfig creates a new RabbitMQConfig.
func NewRabbitMQConfig() *RabbitMQConfig {
	viper.SetConfigName("config")
	viper.AddConfigPath("config") // or the path where your config file is located
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Failed to read config file: ", err)
	}

	return &RabbitMQConfig{
		URL:      viper.GetString("rabbitmq.url"),
		Exchange: viper.GetString("rabbitmq.exchange"),
	}
}

// PublishMessage publishes a message to RabbitMQ.
func (config *RabbitMQConfig) PublishMessage(message string) error {
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		config.Exchange, // exchange name
		"fanout",        // exchange type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return err
	}

	err = ch.Publish(
		config.Exchange, // exchange
		"",              // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	return err
}

func NewRabbitMQPublisher(config *RabbitMQConfig) (*RabbitMQPublisher, error) {
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQPublisher{
		channel:  ch,
		exchange: config.Exchange,
	}, nil
}

func (publisher *RabbitMQPublisher) Publish(message string) error {
	err := publisher.channel.ExchangeDeclare(
		publisher.exchange, // exchange name
		"fanout",           // exchange type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return err
	}

	err = publisher.channel.Publish(
		publisher.exchange, // exchange
		"",                 // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	return err
}

func NewRabbitMQConsumer(config *RabbitMQConfig) (*RabbitMQConsumer, error) {
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	queue, err := ch.QueueDeclare(
		"default-queue", // queue name
		false,           // durable
		false,           // auto-delete
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return nil, err
	}

	messages, err := ch.Consume(
		queue.Name,          // queue
		"news-agg-consumer", // consumer
		true,                // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQConsumer{
		channel:  ch,
		queue:    &queue,
		messages: messages,
	}, nil
}

func (consumer *RabbitMQConsumer) Messages() <-chan amqp.Delivery {
	return consumer.messages
}
