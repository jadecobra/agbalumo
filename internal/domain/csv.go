package domain

// BulkUploadResult tracks the outcome of a bulk upload operation
type BulkUploadResult struct {
	TotalProcessed int      `json:"total_processed"`
	SuccessCount   int      `json:"success_count"`
	FailureCount   int      `json:"failure_count"`
	Errors         []string `json:"errors"`
}
