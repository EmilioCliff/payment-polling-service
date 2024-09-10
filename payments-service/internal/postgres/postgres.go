package postgres

import (
	"context"
	"fmt"
	"sync"

	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/postgres/generated"
	"github.com/EmilioCliff/payment-polling-app/payment-service/pkg"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	storeInstance *Store
	once          sync.Once
)

type Store struct {
	Queries *generated.Queries
	config  pkg.Config
}

func GetStore() (*Store, error) {
	var err error
	once.Do(func() {
		storeInstance, err = NewStore()
	})
	return storeInstance, err
}

func NewStore() (*Store, error) {
	store := &Store{}

	err := store.Start()
	if err != nil {
		return nil, fmt.Errorf("error starting store: %v", err)
	}

	return store, nil
}

func (s *Store) Start() error {
	config, err := pkg.LoadConfig(".")
	if err != nil {
		return err
	}
	s.config = config

	if s.config.DB_URL == "" {
		return fmt.Errorf("dsn is empty")
	}

	postgresConn, err := pgxpool.New(context.Background(), s.config.DB_URL)
	if err != nil {
		return fmt.Errorf("Failed to connect to db: %s", err)
	}

	queries := generated.New(postgresConn)
	s.Queries = queries

	err = s.migrate()
	if err != nil {
		return err
	}

	return nil
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
