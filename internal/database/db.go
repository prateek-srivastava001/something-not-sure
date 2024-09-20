package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
)

var DB *sql.DB

func ConnectPostgres() {
	var err error
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set in the environment")
	}

	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to NeonDB: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Cannot ping NeonDB: %v", err)
	}

	fmt.Println("Successfully connected to NeonDB")
}

func DisconnectPostgres() {
	if DB != nil {
		DB.Close()
	}
}

func RunMigrations() {
	if DB == nil {
		log.Fatal("Database connection is not initialized")
	}

	migrationsDir := "./db/migrations"
	if err := goose.Up(DB, migrationsDir); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	fmt.Println("Migrations applied successfully.")
}
