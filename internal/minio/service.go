package minio

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"hris/internal/user/repository"
)

var (
	ErrUserHasNoProfileImageToDelete = errors.New("user has no profile image to delete")
)

type MinioService interface {
	UploadUserImage(ctx context.Context, userID string, file multipart.File, header *multipart.FileHeader) (string, error)
	UpdateUserImage(ctx context.Context, userID string, file multipart.File, header *multipart.FileHeader) (string, error)
	DeleteUserImage(ctx context.Context, userID string) error
}

type minioService struct {
	minioRepo MinioRepository
	userRepo  repository.UserRepository
}

func NewMinioService(minioRepo MinioRepository, userRepo repository.UserRepository) MinioService {
	return &minioService{
		minioRepo: minioRepo,
		userRepo:  userRepo,
	}
}

// UploadUserImage uploads a new profile image for a user
func (s *minioService) UploadUserImage(ctx context.Context, userID string, file multipart.File, header *multipart.FileHeader) (string, error) {
	// Validate the file
	if err := ValidateImageFile(header); err != nil {
		return "", err
	}

	// Check if user exists
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	// Generate object name with best practice folder structure
	objectName := s.generateObjectName(userID, header.Filename)

	// Get content type
	contentType := header.Header.Get("Content-Type")

	// Upload to MinIO
	fileURL, err := s.minioRepo.Upload(ctx, file, objectName, contentType, header.Size)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Update user's profile image URL in database
	// Update user's profile image URL in database
	user.ProfileImageUrl = fileURL
	// Repo handles updated_at automatically based on SQL implementation, but entity should have it too
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		// Attempt to delete uploaded file if database update fails
		_ = s.minioRepo.Delete(ctx, objectName)
		return "", fmt.Errorf("failed to update user profile: %w", err)
	}

	return fileURL, nil
}

// UpdateUserImage replaces an existing profile image
func (s *minioService) UpdateUserImage(ctx context.Context, userID string, file multipart.File, header *multipart.FileHeader) (string, error) {
	// Validate the new file
	if err := ValidateImageFile(header); err != nil {
		return "", err
	}

	// Get user to check for existing image
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	// Delete old image if exists
	if user.ProfileImageUrl != "" {
		oldObjectName := s.extractObjectNameFromURL(user.ProfileImageUrl)
		if oldObjectName != "" {
			// Ignore error if old file doesn't exist
			_ = s.minioRepo.Delete(ctx, oldObjectName)
		}
	}

	// Generate new object name
	objectName := s.generateObjectName(userID, header.Filename)

	// Get content type
	contentType := header.Header.Get("Content-Type")

	// Upload new image
	fileURL, err := s.minioRepo.Upload(ctx, file, objectName, contentType, header.Size)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Update user's profile image URL
	// Update user's profile image URL
	user.ProfileImageUrl = fileURL
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		// Attempt to delete newly uploaded file if database update fails
		_ = s.minioRepo.Delete(ctx, objectName)
		return "", fmt.Errorf("failed to update user profile: %w", err)
	}

	return fileURL, nil
}

// DeleteUserImage removes a user's profile image
func (s *minioService) DeleteUserImage(ctx context.Context, userID string) error {
	// Get user to check for existing image
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Check if user has a profile image
	if user.ProfileImageUrl == "" {
		return ErrUserHasNoProfileImageToDelete
	}

	// Extract object name from URL
	objectName := s.extractObjectNameFromURL(user.ProfileImageUrl)
	if objectName == "" {
		return errors.New("invalid profile image URL")
	}

	// Delete from MinIO
	if err := s.minioRepo.Delete(ctx, objectName); err != nil {
		return fmt.Errorf("failed to delete file from storage: %w", err)
	}

	// Update user's profile image URL to empty string
	// Update user's profile image URL to empty string
	user.ProfileImageUrl = ""
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user profile: %w", err)
	}

	return nil
}

// generateObjectName creates a storage path following best practices
// Format: users/{userId}/profile/{timestamp}-{randomHash}.{ext}
// Example: users/92/profile/2025-11-28T12:11:45Z-a81fbc2.png
func (s *minioService) generateObjectName(userID, filename string) string {
	// Get file extension (lowercase)
	ext := strings.ToLower(filepath.Ext(filename))
	if len(ext) > 0 && ext[0] == '.' {
		ext = ext[1:]
	}

	// Generate timestamp in ISO 8601 format (UTC)
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05Z")

	// Generate random hash (7 characters for readability)
	randomHash := uuid.New().String()[:7]

	// Build object name: users/{userId}/profile/{timestamp}-{hash}.{ext}
	objectName := fmt.Sprintf("users/%s/profile/%s-%s.%s", userID, timestamp, randomHash, ext)

	return objectName
}

// extractObjectNameFromURL extracts the object name from a full MinIO URL
// Example: http://localhost:9000/mybucket/users/92/profile/2025-11-28T12:11:45Z-a81fbc2.png
// Returns: users/92/profile/2025-11-28T12:11:45Z-a81fbc2.png
func (s *minioService) extractObjectNameFromURL(fileURL string) string {
	// Split URL by "/"
	parts := strings.Split(fileURL, "/")

	// Find the bucket name and extract everything after it
	// URL format: protocol://endpoint/bucket/objectName
	if len(parts) < 5 {
		return ""
	}

	// Join everything after the bucket (index 4 onwards)
	objectName := strings.Join(parts[4:], "/")
	return objectName
}
