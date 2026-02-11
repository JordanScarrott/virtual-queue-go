package db

import (
	"context"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
)

// ConnectAndMigrate connects to the database and runs migrations
func ConnectAndMigrate(connString string) (*pgx.Conn, error) {
	// Run migrations
	m, err := migrate.New("file://db/migrations", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Migrations applied successfully")

	// Connect to database using pgx
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return conn, nil
}
