package kafka

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Consumer struct {
	log     *zap.SugaredLogger
	conn    *kafka.Conn
	topic   string
	groupID string
}

func NewConsumer(log *zap.SugaredLogger, conn *kafka.Conn, topic, groupID string) *Consumer {
	return &Consumer{log: log, conn: conn, topic: topic, groupID: groupID}
}

func (c *Consumer) Consume(ctx context.Context, handler func(ctx context.Context, message *kafka.Message) error) error {
	br := c.conn.Broker()
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{net.JoinHostPort(br.Host, strconv.Itoa(br.Port))},
		Topic:   c.topic,
		GroupID: c.groupID,
	})

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			c.log.Errorf("failed to read message: %v", err)

			break
		}

		if err := handler(ctx, &m); err != nil {
			c.log.Errorf("failed to handle message: %v", err)

			break
		}

		if err := r.CommitMessages(ctx, m); err != nil {
			c.log.Errorf("failed to commit message: %v", err)
		}
	}

	return nil
}

func (c *Consumer) Close() error {
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("failed to close kafka connection: %w", err)
	}

	return nil
}
