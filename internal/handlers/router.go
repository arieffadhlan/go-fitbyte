package handlers

import (
	"github.com/arieffadhlan/go-fitbyte/internal/config"
	activityRepository "github.com/arieffadhlan/go-fitbyte/internal/repositories/activity"
	activityUseCase "github.com/arieffadhlan/go-fitbyte/internal/usecases/activity"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func SetupRouter(cfg *config.Config, db *sqlx.DB, app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"Status": "Welcome to Fitbyte"})
	})

	app.Get("/health-check", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "Ok"})
	})

	version1 := app.Group("v1")

	activityRepository := activityRepository.NewActivityRepository(db)
	activityUseCase := activityUseCase.NewUseCase(activityRepository)
	activityHandler := NewActivityHandler(activityUseCase)

	activityRouter := version1.Group("activity")
	activityRouter.Post("/:user_id", activityHandler.Post)
	activityRouter.Patch("/:id", activityHandler.Update)
}
