package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/arieffadhlan/go-fitbyte/internal/usecases/file"
	"github.com/gofiber/fiber/v2"
)

type FileHandler interface {
	Post(ctx *fiber.Ctx) error
}

type fileHandler struct {
	fileUseCase file.UseCase
}

func NewFileHandler(fileUseCase file.UseCase) FileHandler {
	return &fileHandler{
		fileUseCase: fileUseCase,
	}
}

func (r *fileHandler) Post(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("file")

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	src, err := file.Open()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	defer src.Close()

	if file.Size > (100 * 1024) {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "file exceeds the maximum limit of 100KiB",
		})
	}

	fileName := file.Filename
	fileType := file.Header.Get("Content-Type")

	if !isAllowedFileType(fileName, fileType) {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "file type is not allowed",
		})
	}

	filename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), fileName)
	publicUrl, err := r.fileUseCase.UploadFile(ctx.Context(), file, src, filename)
	if err != nil {
		return ctx.JSON(fiber.Map{"message": err.Error()})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"uri": publicUrl,
	})
}

func isAllowedFileType(fileName, fileType string) bool {
	allowedMimeTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		// "application/octet-stream": true,
	}

	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	if !allowedMimeTypes[fileType] {
		return false
	}

	if fileType == "application/octet-stream" {
		ext := strings.ToLower(filepath.Ext(fileName))
		if !allowedExtensions[ext] {
			return false
		}
	}

	return true
}
