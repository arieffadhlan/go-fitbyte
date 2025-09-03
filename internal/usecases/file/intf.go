package file

import (
	"context"
	"mime/multipart"
)

type UseCase interface {
	UploadFile(context.Context, *multipart.FileHeader, multipart.File, string) (string, error)
}
