package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	fileUseCase "github.com/arieffadhlan/go-fitbyte/internal/usecases/file"
	userUseCase "github.com/arieffadhlan/go-fitbyte/internal/usecases/user"
	"github.com/gofiber/fiber/v2"
)

type FileHandler interface {
	Post(ctx *fiber.Ctx) error
}

type fileHandler struct {
	fileUseCase fileUseCase.UseCase
	userUseCase userUseCase.UseCase
}

func NewFileHandler(fileUseCase fileUseCase.UseCase, userUseCase userUseCase.UseCase) FileHandler {
	return &fileHandler{
		fileUseCase: fileUseCase,
		userUseCase: userUseCase,
	}
}

func (r *fileHandler) Post(ctx *fiber.Ctx) error {
	userId := ctx.Locals("id").(int)
	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	src, err := file.Open()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer src.Close()

	if file.Size > (100 * 1024) {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "file exceeds the maximum limit of 100KiB"})
	}

	fileName := file.Filename
	fileType := file.Header.Get("Content-Type")

	if !isAllowedFileType(fileName, fileType) {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "file type is not allowed",
		})
	}

	user, err := r.userUseCase.GetUserById(ctx.Context(), userId)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if user.ImageURI != nil {
		_, objectName, err := splitPresignedURL(*user.ImageURI)
		if err == nil {
			_ = r.fileUseCase.DeleteFile(ctx.Context(), objectName)
		}
	}

	filename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), fileName)
	imageUri, err := r.fileUseCase.UploadFile(ctx.Context(), file, src, filename)
	if err != nil {
		return ctx.JSON(fiber.Map{"message": err.Error()})
	}

	if err := r.userUseCase.UpdateUserImage(ctx.Context(), userId, imageUri); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"uri": imageUri,
	})
}

func splitPresignedURL(imageURI string) (bucket, objectName string, err error) {
	u, err := url.Parse(imageURI)
	if err != nil {
		return "", "", err
	}

	// Path format: /<bucket>/<objectName>
	parts := strings.SplitN(u.Path, "/", 3)
	if len(parts) < 3 {
		return "", "", fmt.Errorf("invalid URL path: %s", u.Path)
	}

	bucket = parts[1]
	objectName = parts[2]
	return bucket, objectName, nil
}

func isAllowedFileType(fileName, fileType string) bool {
	allowedMimeTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		// "application/octet-stream": true,
	}

	allowedExtensions := map[string]bool{
		".jpeg": true,
		".jpg":  true,
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
