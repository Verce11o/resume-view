package kafka

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/segmentio/kafka-go"
)

const (
	bufferMessageAmount = 1
)

type Notifier struct {
	conn   *kafka.Conn
	writer *kafka.Writer
	topic  string
}

func NewNotifier(conn *kafka.Conn, topic string) *Notifier {
	br := conn.Broker()

	writer := &kafka.Writer{
		Addr:      kafka.TCP(net.JoinHostPort(br.Host, strconv.Itoa(br.Port))),
		Topic:     topic,
		Balancer:  &kafka.LeastBytes{},
		BatchSize: bufferMessageAmount,
	}

	return &Notifier{conn: conn, topic: topic, writer: writer}
}

func (n *Notifier) SendMessage(ctx context.Context, key, value []byte) error {
	err := n.writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})

	if err != nil {
		return fmt.Errorf("could not send message: %w", err)
	}

	return nil
}
