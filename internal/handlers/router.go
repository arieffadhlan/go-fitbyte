package handlers

import (
	"github.com/arieffadhlan/go-fitbyte/internal/config"
	"github.com/arieffadhlan/go-fitbyte/internal/pkg/jwt"
	activityRepository "github.com/arieffadhlan/go-fitbyte/internal/repositories/activity"
	AuthRepository "github.com/arieffadhlan/go-fitbyte/internal/repositories/auth"
	UserRepository "github.com/arieffadhlan/go-fitbyte/internal/repositories/user"
	profileRepository "github.com/arieffadhlan/go-fitbyte/internal/repositories/profile"
	activityUseCase "github.com/arieffadhlan/go-fitbyte/internal/usecases/activity"
	AuthUseCase "github.com/arieffadhlan/go-fitbyte/internal/usecases/auth"
	UserUseCase "github.com/arieffadhlan/go-fitbyte/internal/usecases/user"
	profileUseCase "github.com/arieffadhlan/go-fitbyte/internal/usecases/profile"
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

	handleAuthorization := jwt.Middleware(cfg.JwtSecret)

	activityRepository := activityRepository.NewActivityRepository(db)
	activityUseCase := activityUseCase.NewUseCase(activityRepository)
	activityHandler := NewActivityHandler(activityUseCase)

	activityRouter := v1.Group("activity", handleAuthorization)
	activityRouter.Post("/", activityHandler.Post)
	activityRouter.Patch("/:id", activityHandler.Update)
	activityRouter.Get("/", activityHandler.GetAll)
	activityRouter.Get("/:id", activityHandler.GetById)

	userRepo := UserRepository.NewUserRepository(db)
	userUsecase := UserUseCase.NewUserUsecase(userRepo)

	fileUsecase := FileUseCase.NewUseCase(*cfg)
	fileHandler := NewFileHandler(fileUsecase, userUsecase)

	fileRouter := v1.Group("/file", handleAuthorization)
	fileRouter.Post("/", fileHandler.Post)

	// Profile routes
	profileRepo := profileRepository.NewProfileRepository(db)
	profileUsecase := profileUseCase.NewProfileUseCase(profileRepo)
	profileHandler := NewProfileHandler(profileUsecase)

	profileRouter := v1.Group("/user", handleAuthorization)

	profileRouter.Get("/", profileHandler.GetProfile)
	profileRouter.Patch("/", profileHandler.UpdateProfile)
}
