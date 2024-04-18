package postgres

import (
	"context"
	"fmt"
	"github.com/Verce11o/resume-view/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

func Run(ctx context.Context, cfg *config.Config) *pgxpool.Pool {
	db, err := pgxpool.New(ctx, fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name, cfg.DB.SSLMode))

	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	if err = db.Ping(ctx); err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	return db
}
