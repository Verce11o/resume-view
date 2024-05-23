package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Verce11o/resume-view/employee-service/internal/domain"
	"github.com/Verce11o/resume-view/employee-service/internal/lib/customerrors"
	"github.com/Verce11o/resume-view/employee-service/internal/repository/mongodb"
	"github.com/Verce11o/resume-view/employee-service/internal/repository/postgres"
	pqcache "github.com/Verce11o/resume-view/employee-service/internal/repository/pqCache"
	pqnocache "github.com/Verce11o/resume-view/employee-service/internal/repository/pqNoCache"
	"github.com/Verce11o/resume-view/employee-service/internal/service"
	_ "github.com/flashlabs/rootpath"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type benchRepo interface {
	service.PositionRepository
	terminate(ctx context.Context) error
	benchName(method string) string
}

type bench struct {
	service.PositionRepository
	container testcontainers.Container
	name      string
}

func (b *bench) terminate(ctx context.Context) error {
	err := b.container.Terminate(ctx)
	if err != nil {
		return fmt.Errorf("could not terminate %s container: %w", b.name, err)
	}

	return nil
}

func (b *bench) benchName(method string) string {
	return fmt.Sprintf("%s_%s", b.name, method)
}

func newMongoBench(ctx context.Context, t testing.TB) benchRepo {
	container, connURI := mongodb.SetupMongoContainer(ctx, t)

	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(connURI),
		options.Client().SetMaxConnIdleTime(3*time.Second))
	require.NoError(t, err)

	repo := mongodb.NewPositionRepository(client.Database("employees"))

	return &bench{repo, container, "mongo"}
}

func newPgxBench(ctx context.Context, t testing.TB) benchRepo {
	container, connURI := postgres.SetupPostgresContainer(ctx, t)
	dbPool, err := pgxpool.New(ctx, connURI)
	require.NoError(t, err)

	repo := postgres.NewPositionRepository(dbPool)

	return &bench{repo, container, "postgres"}
}

func newPqNoCacheBench(ctx context.Context, t testing.TB) benchRepo {
	container, connURI := postgres.SetupPostgresContainer(ctx, t)

	db, err := sql.Open("postgres", connURI)
	require.NoError(t, err)

	repo := pqnocache.NewPositionRepository(db)

	return &bench{repo, container, "pq no-cache"}
}

func newPqCacheBench(ctx context.Context, t testing.TB) benchRepo {
	container, connURI := postgres.SetupPostgresContainer(ctx, t)

	db, err := sql.Open("postgres", connURI)
	require.NoError(t, err)

	repo, err := pqcache.NewPositionRepository(db)
	require.NoError(t, err)

	return &bench{repo, container, "pq cache statements"}
}

func BenchmarkPositionRepository(b *testing.B) {
	ctx := context.TODO()

	repos := map[string]benchRepo{
		"postgres":            newPgxBench(ctx, b),
		"mongo":               newMongoBench(ctx, b),
		"pq no-cache":         newPqNoCacheBench(ctx, b),
		"pq cache statements": newPqCacheBench(ctx, b),
	}

	b.ResetTimer()

	for _, repo := range repos {
		defer func(repo benchRepo, ctx context.Context) {
			err := repo.terminate(ctx)
			if err != nil {
				b.Fatalf("failed to terminate benchmark repository: %s", err)
			}
		}(repo, ctx)

		b.Run(repo.benchName("CreatePosition"), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				req := domain.CreatePosition{
					ID:     uuid.New(),
					Name:   "Software Engineer",
					Salary: 60000,
				}
				_, err := repo.CreatePosition(ctx, req)

				if err != nil {
					b.Fatalf("failed to create position: %v", err)
				}
			}
		})

		b.Run(repo.benchName("GetPosition"), func(b *testing.B) {
			req := domain.CreatePosition{
				ID:     uuid.New(),
				Name:   "Software Engineer",
				Salary: 60000,
			}
			position, err := repo.CreatePosition(ctx, req)

			if err != nil {
				b.Fatalf("failed to create position: %v", err)
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, err := repo.GetPosition(ctx, position.ID)
				if err != nil {
					b.Fatalf("failed to get position: %v", err)
				}
			}
		})

		b.Run(repo.benchName("GetPositionList"), func(b *testing.B) {
			for i := 0; i < 10; i++ {
				req := domain.CreatePosition{
					ID:     uuid.New(),
					Name:   "Software Engineer",
					Salary: 60000,
				}
				_, err := repo.CreatePosition(ctx, req)

				if err != nil {
					b.Fatalf("failed to create position: %v", err)
				}
			}

			b.ResetTimer()

			var nextCursor string
			for i := 0; i < b.N; i++ {
				resp, err := repo.GetPositionList(ctx, nextCursor)
				if err != nil {
					b.Fatalf("failed to get position list: %v", err)
				}

				nextCursor = resp.Cursor
			}
		})

		b.Run(repo.benchName("UpdatePosition"), func(b *testing.B) {
			req := domain.CreatePosition{
				ID:     uuid.New(),
				Name:   "Software Engineer",
				Salary: 60000,
			}
			position, err := repo.CreatePosition(ctx, req)

			if err != nil {
				b.Fatalf("failed to create position: %v", err)
			}

			updateReq := domain.UpdatePosition{
				ID:     position.ID,
				Name:   "Senior Software Engineer",
				Salary: 80000,
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, err := repo.UpdatePosition(ctx, updateReq)
				if err != nil {
					b.Fatalf("failed to update position: %v", err)
				}
			}
		})

		b.Run(repo.benchName("DeletePosition"), func(b *testing.B) {
			req := domain.CreatePosition{
				ID:     uuid.New(),
				Name:   "Software Engineer",
				Salary: 60000,
			}
			position, err := repo.CreatePosition(ctx, req)

			if err != nil {
				b.Fatalf("failed to create position: %v", err)
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				err := repo.DeletePosition(ctx, position.ID)
				if err != nil && !errors.Is(err, customerrors.ErrPositionNotFound) {
					b.Fatalf("failed to delete position: %v", err)
				}
			}
		})
	}
}
