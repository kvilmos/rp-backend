package repository

import (
	"context"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
)

type FileStore interface {
	GenerateSignedUploadUrl(ctx context.Context, bucket string, fileKey string, ttl time.Duration) (*url.URL, error)
	GenerateSignedDownloadUrl(ctx context.Context, bucket string, fileKey string, ttl time.Duration) (*url.URL, error)
}

type fileRepository struct {
	minioClient *minio.Client
}

func NewFileRepository(client *minio.Client) FileStore {
	return &fileRepository{
		minioClient: client,
	}
}

func (r *fileRepository) GenerateSignedUploadUrl(ctx context.Context, bucket string, fileKey string, ttl time.Duration) (*url.URL, error) {
	presignedUrl, err := r.minioClient.PresignedPutObject(ctx, bucket, fileKey, ttl)
	return presignedUrl, err
}

func (r *fileRepository) GenerateSignedDownloadUrl(ctx context.Context, bucket string, fileKey string, ttl time.Duration) (*url.URL, error) {
	return r.minioClient.PresignedGetObject(context.Background(), bucket, fileKey, ttl, nil)
}
