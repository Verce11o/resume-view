package nats

import (
	"context"
	"github.com/Verce11o/resume-view/internal/config"
	"github.com/nats-io/nats.go"
	"log"
)

const (
	subject = "views"
)

type Nats struct {
	Conn *nats.Conn
}

func NewNats(ctx context.Context, cfg *config.Config) *Nats {

	nc, err := nats.Connect(cfg.Nats.Url)
	if err != nil {
		log.Fatalf("error connecting to nats: %v", err)
	}

	return &Nats{Conn: nc}

}
