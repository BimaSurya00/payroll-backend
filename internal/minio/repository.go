package minio

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
)

type MinioRepository interface {
	Upload(ctx context.Context, file io.Reader, objectName, contentType string, fileSize int64) (string, error)
	Delete(ctx context.Context, objectName string) error
	ObjectExists(ctx context.Context, objectName string) (bool, error)
}

type minioRepository struct {
	client   *MinioClient
	endpoint string
}

func NewMinioRepository(client *MinioClient, endpoint string) MinioRepository {
	return &minioRepository{
		client:   client,
		endpoint: endpoint,
	}
}

// Upload uploads a file to MinIO and returns the public URL
func (r *minioRepository) Upload(ctx context.Context, file io.Reader, objectName, contentType string, fileSize int64) (string, error) {
	// Upload object with content type
	_, err := r.client.Client.PutObject(
		ctx,
		r.client.Bucket,
		objectName,
		file,
		fileSize,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to MinIO: %w", err)
	}

	// Return public URL
	publicURL := r.client.GetPublicURL(r.endpoint, objectName)
	return publicURL, nil
}

// Delete removes an object from MinIO
func (r *minioRepository) Delete(ctx context.Context, objectName string) error {
	// Check if object exists first
	exists, err := r.ObjectExists(ctx, objectName)
	if err != nil {
		return fmt.Errorf("failed to check object existence: %w", err)
	}

	if !exists {
		return fmt.Errorf("object not found: %s", objectName)
	}

	// Remove object
	err = r.client.Client.RemoveObject(ctx, r.client.Bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object from MinIO: %w", err)
	}

	return nil
}

// ObjectExists checks if an object exists in MinIO
func (r *minioRepository) ObjectExists(ctx context.Context, objectName string) (bool, error) {
	_, err := r.client.Client.StatObject(ctx, r.client.Bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		// Check if error is "object not found"
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}