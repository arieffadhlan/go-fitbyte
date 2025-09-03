package handlers

import (
	"net/http"

	"github.com/arieffadhlan/go-fitbyte/internal/dto"
	"github.com/arieffadhlan/go-fitbyte/internal/usecases/auth"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	uc *auth.AuthUsecase
}

func NewAuthHandler(uc *auth.AuthUsecase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func (ah *AuthHandler) Register(ctx *fiber.Ctx) error {
	var req dto.AuthRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request payload",
		})
	}

	id, err := ah.uc.Register(ctx.Context(), &req)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"id":    id,
		"email": req.Email,
	})
}

func (ah *AuthHandler) Login(ctx *fiber.Ctx) error {
	var req dto.AuthRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request payload",
		})
	}

	res, err := ah.uc.Login(ctx.Context(), &req)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(res)
}
