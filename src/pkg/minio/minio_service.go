package minIOService

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var service MinIOService

type MinIOService struct {
	client *minio.Client
}

func InitMinIO() error {
	endpoint := "localhost:9000"
	accessKeyID := "minioadmin"
	secretAccessKey := "minioadmin"
	useSSL := false

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return err
	}
	service = MinIOService{client: client}

	return nil
}

func GetClient() *minio.Client {
	return service.client
}
