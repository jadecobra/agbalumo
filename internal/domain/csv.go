package domain

import (
	"context"
	"io"
)

// BulkUploadResult tracks the outcome of a bulk upload operation
type BulkUploadResult struct {
	TotalProcessed int      `json:"total_processed"`
	SuccessCount   int      `json:"success_count"`
	FailureCount   int      `json:"failure_count"`
	Errors         []string `json:"errors"`
}

// CSVService defines the contract for processing bulk CSV uploads and exports
type CSVService interface {
	ParseAndImport(ctx context.Context, r io.Reader, repo ListingStore) (*BulkUploadResult, error)
	GenerateCSV(ctx context.Context, listings []Listing) (io.Reader, error)
}
