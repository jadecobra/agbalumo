package domain

import (
	"context"
	"mime/multipart"
)

// ImageService defines the contract for handling image uploads
type ImageService interface {
	UploadImage(ctx context.Context, file *multipart.FileHeader, listingID string) (string, error)
	DeleteImage(ctx context.Context, imageURL string) error
}
