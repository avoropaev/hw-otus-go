package rmq

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"
)

type Producer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	uri     string
	queue   string
}

func NewProducer(uri, queue string) *Producer {
	return &Producer{
		uri:   uri,
		queue: queue,
	}
}

func (c *Producer) Publish(_ context.Context, body []byte) error {
	if c.conn == nil || c.conn.IsClosed() {
		if err := c.Connect(); err != nil {
			return err
		}
	}

	queue, err := c.channel.QueueDeclarePassive(c.queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = c.channel.Publish("", queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         body,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Producer) Connect() error {
	var err error

	c.conn, err = amqp.Dial(c.uri)
	if err != nil {
		return fmt.Errorf("dial: %s", err)
	}

	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("channel: %s", err)
	}

	return nil
}

func (c *Producer) Disconnect() error {
	if c.conn == nil {
		return nil
	}

	return c.conn.Close()
}
