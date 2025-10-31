package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService struct {
	cld *cloudinary.Cloudinary
}

// NewCloudinaryService creates a new Cloudinary service instance
func NewCloudinaryService(cloudinaryURL string) (*CloudinaryService, error) {
	if cloudinaryURL == "" {
		return nil, fmt.Errorf("cloudinary URL is required")
	}

	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Cloudinary: %w", err)
	}

	return &CloudinaryService{
		cld: cld,
	}, nil
}

// UploadImage uploads an image to Cloudinary and returns the URL
func (s *CloudinaryService) UploadImage(ctx context.Context, file multipart.File, filename string) (string, error) {
	// Create a unique public ID for the image
	ext := strings.ToLower(filepath.Ext(filename))
	publicID := fmt.Sprintf("health-store/products/%d_%s", time.Now().Unix(), strings.TrimSuffix(filename, ext))

	// Upload the file to Cloudinary
	uploadResult, err := s.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID:       publicID,
		Folder:         "health-store/products",
		ResourceType:   "image",
		Transformation: "c_limit,w_1000,h_1000,q_auto,f_auto", // Auto optimize image
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload image to Cloudinary: %w", err)
	}

	return uploadResult.SecureURL, nil
}

// DeleteImage deletes an image from Cloudinary
func (s *CloudinaryService) DeleteImage(ctx context.Context, imageURL string) error {
	// Extract public ID from the URL
	publicID := extractPublicIDFromURL(imageURL)
	if publicID == "" {
		return fmt.Errorf("invalid image URL: cannot extract public ID")
	}

	// Delete the image from Cloudinary
	_, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: "image",
	})

	if err != nil {
		return fmt.Errorf("failed to delete image from Cloudinary: %w", err)
	}

	return nil
}

// extractPublicIDFromURL extracts the public ID from a Cloudinary URL
// Example: https://res.cloudinary.com/demo/image/upload/v1234567890/health-store/products/image.jpg
// Returns: health-store/products/image
func extractPublicIDFromURL(url string) string {
	if url == "" {
		return ""
	}

	// Split by "/"
	parts := strings.Split(url, "/")

	// Find the "upload" part and get everything after the version number
	for i, part := range parts {
		if part == "upload" && i+2 < len(parts) {
			// Skip the version (v1234567890) and get the rest
			publicIDParts := parts[i+2:]
			publicID := strings.Join(publicIDParts, "/")

			// Remove file extension
			if lastDot := strings.LastIndex(publicID, "."); lastDot != -1 {
				publicID = publicID[:lastDot]
			}

			return publicID
		}
	}

	return ""
}
