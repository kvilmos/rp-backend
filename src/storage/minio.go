package storage

import (
	"context"
	"os"
	"room-planner/common/constant"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinioClient() (*minio.Client, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_USERNAME")
	secretAccessKey := os.Getenv("MINIO_PASSWORD")
	useSSL := false

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	err = createBuckets(client)
	if err != nil {
		return nil, err

	}

	return client, nil
}

func createBuckets(client *minio.Client) error {
	err := client.MakeBucket(context.Background(), constant.FURNITURE_BUCKET, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := client.BucketExists(context.Background(), constant.FURNITURE_BUCKET)
		if errBucketExists != nil || !exists {
			return err
		}
	}

	err = client.MakeBucket(context.Background(), constant.THUMBNAIL_BUCKET, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := client.BucketExists(context.Background(), constant.THUMBNAIL_BUCKET)
		if errBucketExists != nil || !exists {
			return err
		}
	}

	return nil
}
