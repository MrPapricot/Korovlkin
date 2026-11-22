package main

import (
	"context"
	"fmt"
	"log"
	"main/src/config"
	"main/src/database"
	"main/src/handlers"
	"main/src/repositories"
	"time"

	"github.com/gofiber/fiber/v3"
)

func main() {
	cfg := config.LoadConfig()

	ctx, cancel := context.WithTimeout(
		context.Background(),
		15*time.Second,
	)
	defer cancel()

	pool, err := database.Connect(ctx, cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer pool.Close()
	app := fiber.New()

	taskHandler := handlers.NewTaskHandler(
		repositories.NewTaskRepository(pool),
	)

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	api := app.Group("/api/tasks")

	api.Post("/", taskHandler.Create)
	api.Get("/", taskHandler.List)
	api.Get("/:id", taskHandler.Get)
	api.Put("/:id", taskHandler.Update)
	api.Patch("/:id/status", taskHandler.UpdateStatus)
	api.Delete("/:id", taskHandler.Delete)

	if err := app.Listen(fmt.Sprintf("0.0.0.0:%d", cfg.AppPort)); err != nil {
		log.Fatal(err)
	}
}
