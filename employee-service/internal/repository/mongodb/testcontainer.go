package mongodb

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetupMongoContainer(ctx context.Context, t testing.TB) (testcontainers.Container, string) {
	entrypointScript := []string{
		"/bin/bash", "-c",
		`echo "rs.initiate()" > /docker-entrypoint-initdb.d/1-init-replicaset.js &&
		exec /usr/local/bin/docker-entrypoint.sh mongod --replSet rs0 --bind_ip_all --noauth`,
	}

	mongoContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo:latest",
			ExposedPorts: []string{"27017/tcp"},
			Env: map[string]string{
				"MONGO_APP_DATABASE": "employees",
				"MONGO_REPLICA_PORT": "27018",
			},
			WaitingFor: wait.ForListeningPort("27017/tcp"),
			Entrypoint: entrypointScript,
		},
		Started: true,
	})
	require.NoError(t, err)

	connURI, err := mongoContainer.Endpoint(ctx, "mongodb")
	require.NoError(t, err)

	return mongoContainer, connURI + "/?directConnection=true&tls=false"
}
