package main

import (
	"go-fiber-learning/database"
	"go-fiber-learning/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	//Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	//Connect to database
	database.Connect()

	//Init fiber
	app := fiber.New(fiber.Config{
		AppName: "Go Fiber Learning",
	})

	api := app.Group("/api", logger.New())

	routes.SetupRoutes(api)

	app.Listen(":3000")
}
