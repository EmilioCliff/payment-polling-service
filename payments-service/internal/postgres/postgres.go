package postgres

import (
	"context"
	"fmt"

	"github.com/EmilioCliff/payment-polling-app/payment-service/pkg"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	conn   *pgxpool.Pool
	config pkg.Config
}

func NewStore(config pkg.Config) *Store {
	return &Store{
		config: config,
	}
}

func (s *Store) Start() error {
	if s.config.DB_URL == "" {
		return fmt.Errorf("dsn is empty")
	}

	postgresConn, err := pgxpool.New(context.Background(), s.config.DB_URL)
	if err != nil {
		return fmt.Errorf("Failed to connect to db: %s", err)
	}

	err = postgresConn.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to connect to db")
	}

	s.conn = postgresConn

	return s.migrate()
}

func (s *Store) migrate() error {
	if s.config.MIGRATION_PATH == "" {
		return fmt.Errorf("migration dir is empty")
	}

	migration, err := migrate.New(s.config.MIGRATION_PATH, s.config.DB_URL)
	if err != nil {
		return fmt.Errorf("Failed to load migration: %s", err)
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("Failed to run migrate up: %s", err)
	}

	return nil
}
