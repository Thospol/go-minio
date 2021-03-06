package minio

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

var (
	minioClient = &minio.Client{}
	ctx         = context.Background()
)

// Client minio client interface
type Client interface {
	UploadImage(bucketName string, objectName, filePath string) error
}

// Configuration config minio for new connection
type Configuration struct {
	Host            string
	AccessKeyID     string
	SecretAccessKey string
}

// NewConnection new ftp connection
func NewConnection(config Configuration) (err error) {
	// Initialize minio client object.
	minioClient, err = minio.New(config.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		return err
	}

	return nil
}

type client struct {
	client *minio.Client
}

// GetClient get minio client
func GetClient() Client {
	return &client{
		client: minioClient,
	}
}

// UploadImage upload image
func (m *client) UploadImage(bucketName string, objectName, filePath string) error {
	err := m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "us-east-1"})
	if err != nil {
		exists, errBucketExists := m.client.BucketExists(ctx, bucketName)
		if errBucketExists != nil {
			logrus.Errorf("[UploadImage] check bucket exists error: %s", err)
			return err
		}

		if !exists {
			logrus.Errorf("[UploadImage] make bucket error: %s", err)
			return err
		}
	}

	_, err = m.client.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{})
	if err != nil {
		logrus.Errorf("[UploadImage] put object error: %s", err)
		return err
	}

	return nil
}
