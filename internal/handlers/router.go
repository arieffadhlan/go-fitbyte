package handlers

import (
	"github.com/arieffadhlan/go-fitbyte/internal/config"
	"github.com/arieffadhlan/go-fitbyte/internal/pkg/jwt"
	activityRepository "github.com/arieffadhlan/go-fitbyte/internal/repositories/activity"
	AuthRepository "github.com/arieffadhlan/go-fitbyte/internal/repositories/auth"
	activityUseCase "github.com/arieffadhlan/go-fitbyte/internal/usecases/activity"
	AuthUseCase "github.com/arieffadhlan/go-fitbyte/internal/usecases/auth"
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

	authRepo := AuthRepository.NewAuthRepository(db)
	authUsecase := AuthUseCase.NewAuthUsecase(authRepo, cfg)
	authHandler := NewAuthHandler(authUsecase)

	authRouter := v1.Group("")
	authRouter.Post("/login", authHandler.Login)
	authRouter.Post("/register", authHandler.Register)

	activityRepository := activityRepository.NewActivityRepository(db)
	activityUseCase := activityUseCase.NewUseCase(activityRepository)
	activityHandler := NewActivityHandler(activityUseCase)

	activityRouter := v1.Group("activity")
	activityRouter.Post("/:user_id", activityHandler.Post)
	activityRouter.Patch("/:id", activityHandler.Update)
	activityRouter.Get("/", activityHandler.GetAll)
	activityRouter.Get("/:id", activityHandler.GetById)

	// Test
	v1.Get("/test", jwt.Middleware(cfg.JwtSecret), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"userID": c.Locals("userID"),
			"email":  c.Locals("email"),
		})
	})

	fileUsecase := FileUseCase.NewUseCase(*cfg)
	fileHandler := NewFileHandler(fileUsecase)

	fileRouter := v1.Group("/file")
	fileRouter.Post("/", fileHandler.Post)
}
