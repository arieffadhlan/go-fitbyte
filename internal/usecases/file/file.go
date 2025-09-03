package file

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/arieffadhlan/go-fitbyte/internal/config"
	minioUploader "github.com/arieffadhlan/go-fitbyte/internal/pkg/minio"
	"github.com/minio/minio-go/v7"
)

type useCase struct {
	config config.Config
	minio  *minio.Client
}

func NewUseCase(config config.Config) UseCase {
	minioConfig := &minioUploader.MinioConfig{
		AccessKeyID:     config.Minio.AccessKeyID,
		SecretAccessKey: config.Minio.SecretAccessKey,
		UseSSL:          config.Minio.UseSSL,
		Endpoint:        config.Minio.Endpoint,
	}

	minioClient, _ := minioUploader.NewUploader(minioConfig)

	return &useCase{
		config: config,
		minio:  minioClient,
	}
}

func (uc *useCase) UploadFile(ctx context.Context, file *multipart.FileHeader, src multipart.File, fileName string) (string, error) {
	bucketName := uc.config.Minio.BucketName

	_, err := uc.minio.PutObject(
		ctx,
		bucketName,
		fileName,
		src,
		file.Size,
		minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type")},
	)

	if err != nil {
		return "", err
	}

	presignedURL, err := uc.minio.PresignedGetObject(ctx, bucketName, fileName, time.Hour*24*7, nil)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}
