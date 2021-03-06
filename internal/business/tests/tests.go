// Package tests contains supporting code for running tests.
package tests

import (
	"fmt"
	"math/rand"
	"net"
	"runtime"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/johanronkko/quote-service/internal/business/data/schema"
	"github.com/johanronkko/quote-service/internal/foundation/database"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

// NewUnit creates a test database inside a Docker container. It creates the
// required table structure but the database is otherwise empty.
func NewUnit(tb testing.TB) *sqlx.DB {
	pool, err := dockertest.NewPool("")
	if err != nil {
		tb.Fatalf("Couldn't connect to docker: %s", err)
	}
	pool.MaxWait = 10 * time.Second

	cfg := database.Config{
		User:       "postgres",
		Password:   "postgres",
		Name:       "testdb",
		DisableTLS: true,
	}
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13.2-alpine",
		Env: []string{
			fmt.Sprintf("POSTGRES_USER=%s", cfg.User),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", cfg.Password),
			fmt.Sprintf("POSTGRES_DB=%s", cfg.Name),
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		tb.Fatalf("Couldn't start resource: %s", err)
	}

	resource.Expire(60)

	tb.Cleanup(func() {
		if err := pool.Purge(resource); err != nil {
			tb.Fatalf("Couldn't purge container: %v", err)
		}
	})

	cfg.Host = fmt.Sprintf("%s:5432", resource.Container.NetworkSettings.IPAddress)
	if runtime.GOOS == "darwin" { // MacOS-specific
		cfg.Host = net.JoinHostPort(resource.GetBoundIP("5432/tcp"), resource.GetPort("5432/tcp"))
	}

	db, err := database.Open(cfg)
	if err != nil {
		tb.Fatalf("Opening database connection: %s", err)
	}

	tb.Cleanup(func() {
		if err := db.Close(); err != nil {
			tb.Fatalf("Closing database connection %s", err)
		}
	})

	if err := pool.Retry(func() (err error) {
		return db.Ping()
	}); err != nil {
		tb.Fatalf("Database never ready: %s", err)
	}

	if err := schema.Migrate(db); err != nil {
		tb.Fatalf("Migrating error: %s", err)
	}

	return db
}

// NewIntegration creates a test database inside a Docker container. It creates the
// required table structure and seeds the database with 3 quotes.
func NewIntegration(tb testing.TB) *sqlx.DB {
	db := NewUnit(tb)
	if err := schema.Seed(db); err != nil {
		tb.Fatal(err)
	}
	return db
}

// GenRandomAlpha generates a string with n random alpha characters.
func GenRandomAlpha(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
