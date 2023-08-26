package broker

import (
	"context"
	"errors"
	"fmt"
	"github.com/mjedari/sternx-project/distributor/app/configs"
	"github.com/mjedari/sternx-project/distributor/infra/utils"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"log"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Config  configs.Rabbit
}

func NewRabbitMQ(config configs.Rabbit) (*RabbitMQ, error) {
	ctx := context.TODO()
	healerConfig := configs.Config.GetHealerConfig()
	rabbitRetry, err := utils.Retry(func(ctx context.Context) (any, error) {
		conn, err := amqp.Dial(config.GetURL())
		if err != nil {
			return nil, err
		}

		return conn, nil
	}, healerConfig.GetRetryTimes(), healerConfig.GetRetryDelay())(ctx)

	if err != nil {
		return nil, err
	}
	// here we convert interface datatype to redis.Client
	conn := rabbitRetry.(*amqp.Connection)

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = ch.ExchangeDeclare(config.ExchangeName, "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
		return nil, err
	}

	return &RabbitMQ{Conn: conn, Channel: ch, Config: config}, nil
}

func (r *RabbitMQ) CheckHealth(ctx context.Context) error {
	if r.Conn.IsClosed() {
		return errors.New("connection lost error")
	}
	return nil
}

func (r *RabbitMQ) ResetConnection(ctx context.Context) error {
	newConn, err := NewRabbitMQ(r.Config)
	if err != nil {
		logrus.Errorf("error on re creating rabbit connection: %v", err)
		return err
	}
	// todo: check and fix this
	r.Conn = newConn.Conn

	return nil
}

func (r *RabbitMQ) ProduceOnQueue(ctx context.Context, queue string, message []byte) error {
	err := r.Channel.Publish("", queue, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		})

	return err
}

func (r *RabbitMQ) Consume(ctx context.Context, queueName, key string, apply func(ctx context.Context, message []byte) error) error {

	q, err := r.Channel.QueueDeclare(r.Config.QueueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = r.Channel.QueueBind(q.Name, r.Config.RoutingKey, r.Config.ExchangeName, false, nil)

	msgs, err := r.Channel.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-msgs:
			for msg := range msgs {
				processErr := apply(ctx, msg.Body)

				if processErr == nil {
					msg.Ack(false)
				}
			}
		}
	}
}

func (r *RabbitMQ) Close() error {
	err := r.Conn.Close()
	if err != nil {
		return err
	}

	fmt.Println("Closing rabbitmq...")
	return nil
}
