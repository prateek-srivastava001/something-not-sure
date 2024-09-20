package main

import (
	"EasySplit/internal/database"
	"EasySplit/internal/routers"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	database.ConnectPostgres()
	defer database.DisconnectPostgres()

	database.RunMigrations()

	e := echo.New()
	routers.SetupRoutes(e)
	routers.SetupFriendRoutes(e)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := e.Start(":" + port); err != nil {
		e.Logger.Fatal(err)
	}
}
