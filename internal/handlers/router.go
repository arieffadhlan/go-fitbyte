package handlers

import (
	"github.com/arieffadhlan/go-fitbyte/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	FileUseCase "github.com/arieffadhlan/go-fitbyte/internal/usecases/file"
)

func SetupRouter(cfg *config.Config, db *sqlx.DB, app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"Status": "Welcome to Fitbyte"})
	})

	app.Get("/health-check", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "Ok"})
	})

	v1 := app.Group("/api/v1")

	fileUsecase := FileUseCase.NewUseCase(*cfg)
	fileHandler := NewFileHandler(fileUsecase)

	fileRouter := v1.Group("/file")
	fileRouter.Post("/", fileHandler.Post)
}
