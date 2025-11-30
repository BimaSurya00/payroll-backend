package minio

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

type MinioClient struct {
	Client *minio.Client
	Bucket string
	UseSSL bool
}

// NewMinioClient creates a new MinIO client instance
func NewMinioClient(cfg MinioConfig) (*MinioClient, error) {
	// Initialize MinIO client
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	// Create bucket if it doesn't exist
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
		log.Printf("✅ MinIO bucket '%s' created successfully", cfg.Bucket)
	}

	// Set bucket policy to public read (for images)
	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/*"]
			}
		]
	}`, cfg.Bucket)

	err = client.SetBucketPolicy(ctx, cfg.Bucket, policy)
	if err != nil {
		log.Printf("⚠️  Warning: Failed to set bucket policy: %v", err)
	}

	log.Printf("✅ Connected to MinIO at %s", cfg.Endpoint)

	return &MinioClient{
		Client: client,
		Bucket: cfg.Bucket,
		UseSSL: cfg.UseSSL,
	}, nil
}

// GetPublicURL returns the public URL for an object
func (mc *MinioClient) GetPublicURL(endpoint, objectName string) string {
	protocol := "http"
	if mc.UseSSL {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", protocol, endpoint, mc.Bucket, objectName)
}