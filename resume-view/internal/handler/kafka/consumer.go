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
	log       *zap.SugaredLogger
	conn      *kafka.Conn
	topic     string
	partition int
}

func NewConsumer(log *zap.SugaredLogger, conn *kafka.Conn, topic string, partition int) *Consumer {
	return &Consumer{log: log, conn: conn, topic: topic, partition: partition}
}

func (c *Consumer) Consume(ctx context.Context) error {
	br := c.conn.Broker()
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{net.JoinHostPort(br.Host, strconv.Itoa(br.Port))},
		Topic:     c.topic,
		Partition: c.partition,
	})

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			c.log.Errorf("failed to read message: %v", err)

			break
		}

		c.log.Debugf("message at offset %d: %s", m.Offset, string(m.Value))
	}

	return nil
}

func (c *Consumer) Close() error {
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("failed to close kafka connection: %w", err)
	}

	return nil
}
