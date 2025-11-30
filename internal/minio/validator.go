package minio

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
)

const (
	// MaxImageSize is 5MB
	MaxImageSize = 5 * 1024 * 1024 // 5MB in bytes
)

var (
	// AllowedImageMimeTypes defines permitted MIME types for images
	AllowedImageMimeTypes = map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/webp": true,
	}

	// AllowedImageExtensions defines permitted file extensions
	AllowedImageExtensions = map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}
)

// ValidationError represents a file validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidateImageFile validates an uploaded image file
func ValidateImageFile(header *multipart.FileHeader) error {
	// Validate file size
	if header.Size > MaxImageSize {
		return &ValidationError{
			Field:   "file",
			Message: fmt.Sprintf("file size exceeds maximum allowed size of %d bytes (%.2f MB)", MaxImageSize, float64(MaxImageSize)/(1024*1024)),
		}
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !AllowedImageExtensions[ext] {
		return &ValidationError{
			Field:   "file",
			Message: fmt.Sprintf("invalid file extension '%s'. Allowed extensions: .jpg, .jpeg, .png, .webp", ext),
		}
	}

	// Validate MIME type
	contentType := header.Header.Get("Content-Type")
	if !AllowedImageMimeTypes[contentType] {
		return &ValidationError{
			Field:   "file",
			Message: fmt.Sprintf("invalid content type '%s'. Allowed types: image/jpeg, image/png, image/webp", contentType),
		}
	}

	// Validate filename is not empty
	if header.Filename == "" {
		return &ValidationError{
			Field:   "file",
			Message: "filename cannot be empty",
		}
	}

	return nil
}

// GetFileExtension returns the lowercase file extension without the dot
func GetFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	if len(ext) > 0 && ext[0] == '.' {
		ext = ext[1:]
	}
	return strings.ToLower(ext)
}

// IsImageFile checks if a file is an image based on content type
func IsImageFile(contentType string) bool {
	return AllowedImageMimeTypes[contentType]
}

// ValidateFile is a generic file validator that can be extended
func ValidateFile(header *multipart.FileHeader, maxSize int64, allowedTypes map[string]bool, allowedExts map[string]bool) error {
	if header.Size > maxSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %d bytes", maxSize)
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedExts[ext] {
		return fmt.Errorf("invalid file extension: %s", ext)
	}

	contentType := header.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return fmt.Errorf("invalid content type: %s", contentType)
	}

	return nil
}
