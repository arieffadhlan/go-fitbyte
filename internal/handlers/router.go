package handlers

import (
	"github.com/arieffadhlan/go-fitbyte/internal/config"
	"github.com/arieffadhlan/go-fitbyte/internal/pkg/jwt"
	AuthRepository "github.com/arieffadhlan/go-fitbyte/internal/repositories/auth"
	AuthUseCase "github.com/arieffadhlan/go-fitbyte/internal/usecases/auth"
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

	v1 := app.Group("/api/v1")

	authRepo := AuthRepository.NewAuthRepository(db)
	authUsecase := AuthUseCase.NewAuthUsecase(authRepo, cfg)
	authHandler := NewAuthHandler(authUsecase)

	authRouter := v1.Group("")
	authRouter.Post("/login", authHandler.Login)
	authRouter.Post("/register", authHandler.Register)

	// Test
	v1.Get("/test", jwt.Middleware(cfg.JwtSecret), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"userID": c.Locals("userID"),
			"email":  c.Locals("email"),
		})
	})
}
