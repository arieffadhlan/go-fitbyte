package handlers

import (
	"strconv"

	"github.com/arieffadhlan/go-fitbyte/internal/dto"
	"github.com/arieffadhlan/go-fitbyte/internal/pkg/exceptions"
	internal_validator "github.com/arieffadhlan/go-fitbyte/internal/pkg/validator"
	"github.com/arieffadhlan/go-fitbyte/internal/usecases/activity"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ActivityHandler interface {
	Post(fibCtx *fiber.Ctx) error
	Update(fibCtx *fiber.Ctx) error
	GetAll(fibCtx *fiber.Ctx) error
	GetById(fibCtx *fiber.Ctx) error
}

type activityHandler struct {
	activityUseCase activity.UseCase
}

func NewActivityHandler(activityUseCase activity.UseCase) ActivityHandler {
	return &activityHandler{activityUseCase}
}

func (r *activityHandler) GetAll(fibCtx *fiber.Ctx) error {
	activities, err := r.activityUseCase.GetAllActivities(fibCtx.Context())
	if err != nil {
		return fibCtx.Status(exceptions.MapToHttpStatusCode(err)).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return fibCtx.Status(fiber.StatusOK).JSON(activities)
}

func (r *activityHandler) GetById(fibCtx *fiber.Ctx) error {
	activityId := fibCtx.Params("id")
	id, err := strconv.Atoi(activityId)
	if err != nil {
		return fibCtx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "activityId must be an integer",
		})
	}

	activity, err := r.activityUseCase.GetActivityById(fibCtx.Context(), id)
	if err != nil {
		return fibCtx.Status(exceptions.MapToHttpStatusCode(err)).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return fibCtx.Status(fiber.StatusOK).JSON(activity)
}

func (r *activityHandler) Post(fibCtx *fiber.Ctx) error {
	if fibCtx.Get("Content-Type") != "application/json" {
		return fibCtx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid content type",
		})
	}

	var activityRequest dto.ActivityRequest
	if err := fibCtx.BodyParser(&activityRequest); err != nil {
		return fibCtx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var validate = validator.New()
	validate.RegisterValidation("iso8601", internal_validator.ValidateISODate)
	if err := validate.Struct(activityRequest); err != nil {
		return fibCtx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userId, ok := fibCtx.Locals("id").(int)
	if !ok {
		return fibCtx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user id cant be empty",
		})
	}

	activityResponse, err := r.activityUseCase.PostActivity(fibCtx.Context(), &activityRequest, userId)

	if err != nil {
		return fibCtx.Status(exceptions.MapToHttpStatusCode(err)).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return fibCtx.Status(fiber.StatusOK).JSON(activityResponse)
}

func (r *activityHandler) Update(fibCtx *fiber.Ctx) error {
	if fibCtx.Get("Content-Type") != "application/json" {
		return fibCtx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid content type",
		})
	}

	activityId := fibCtx.Params("id")
	if activityId == "" {
		return fibCtx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "activityId is required",
		})
	}

	var activityRequest dto.ActivityRequest
	if err := fibCtx.BodyParser(&activityRequest); err != nil {
		return fibCtx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var validate = validator.New()
	validate.RegisterValidation("iso8601", internal_validator.ValidateISODate)
	if err := validate.Struct(activityRequest); err != nil {
		return fibCtx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userId, ok := fibCtx.Locals("id").(int)
	if !ok {
		return fibCtx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user id cant be empty",
		})
	}

	activityResponse, err := r.activityUseCase.UpdateActivity(fibCtx.Context(), &activityRequest, userId, activityId)

	if err != nil {
		if err == exceptions.ErrNotFound {
			return fibCtx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "activityId is not found",
			})
		}
		return fibCtx.Status(exceptions.MapToHttpStatusCode(err)).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return fibCtx.Status(fiber.StatusOK).JSON(activityResponse)
}
