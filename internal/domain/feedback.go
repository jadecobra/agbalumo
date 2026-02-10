package domain

import "time"

type FeedbackType string

const (
	FeedbackTypeIssue   FeedbackType = "Issue"
	FeedbackTypeFeature FeedbackType = "Feature"
	FeedbackTypeOther   FeedbackType = "Other"
)

type Feedback struct {
	ID        string       `json:"id"`
	UserID    string       `json:"user_id"`
	Type      FeedbackType `json:"type"`
	Content   string       `json:"content"`
	CreatedAt time.Time    `json:"created_at"`
}
