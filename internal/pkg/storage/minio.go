package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOStorage struct {
	client     *minio.Client
	bucketName string
}

func NewMinIOStorage(endpoint, accessKeyID, secretAccessKey, bucketName string, useSSL bool) (*MinIOStorage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}

	return &MinIOStorage{
		client:     client,
		bucketName: bucketName,
	}, nil
}

func (s *MinIOStorage) Upload(ctx context.Context, objectName string, reader io.Reader, objectSize int64, contentType string) (string, error) {
	_, err := s.client.PutObject(ctx, s.bucketName, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	return s.GetURL(objectName), nil
}

func (s *MinIOStorage) UploadFile(ctx context.Context, filePath, objectName, contentType string) (string, error) {
	_, err := s.client.FPutObject(ctx, s.bucketName, objectName, filePath, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	return s.GetURL(objectName), nil
}

func (s *MinIOStorage) Delete(ctx context.Context, objectName string) error {
	return s.client.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{})
}

func (s *MinIOStorage) GetURL(objectName string) string {
	return fmt.Sprintf("/%s/%s", s.bucketName, objectName)
}

func (s *MinIOStorage) GetPresignedURL(ctx context.Context, objectName string, expires time.Duration) (string, error) {
	return s.client.PresignedGetObject(ctx, s.bucketName, objectName, expires, nil)
}

func (s *MinIOStorage) GenerateObjectName(resourceType, extension string) string {
	now := time.Now()
	timestamp := now.Format("20060102150405")
	return filepath.Join(resourceType, now.Format("2006/01"), fmt.Sprintf("%s%s", timestamp, extension))
}
