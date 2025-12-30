package database

import (
	"context"
	"fmt"
	"log"
	"project-app-inventory-restapi-golang-anas/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(cfg *config.Config) *pgxpool.Pool {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Unable to parse database config: %v", err)
	}

	// Optional: Setting connection pool settings
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour

	conn, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	// Test connection (Ping)
	if err := conn.Ping(context.Background()); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	log.Println("Connected to PostgreSQL successfully!")
	return conn
}
