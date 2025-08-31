package app

import (
	"time"

	"github.com/arieffadhlan/go-fitbyte/internal/config"
	"github.com/arieffadhlan/go-fitbyte/internal/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func NewServer(cfg *config.Config, db *sqlx.DB) *fiber.App {
	app := fiber.New(fiber.Config{
		IdleTimeout:  600 * time.Second,
		ReadTimeout:  600 * time.Second,
		WriteTimeout: 600 * time.Second,
	})

	handlers.SetupRouter(cfg, db, app)

	return app
}
