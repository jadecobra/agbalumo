package main

import (
	"fmt"
	"os"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/maintenance"
)

var feedbackListCmd = makeSimpleCmd("feedback-list", "List all feedback submissions locally", func() error {
	dbPath := os.Getenv("DATABASE_URL")
	if dbPath == "" {
		dbPath = domain.DefaultDatabaseURL
	}

	feedbacks, err := maintenance.CheckFeedbackList(dbPath)
	if err != nil {
		return err
	}

	if len(feedbacks) == 0 {
		fmt.Println("No feedback submissions found.")
		return nil
	}

	fmt.Printf("Found %d feedback submission(s):\n\n", len(feedbacks))
	for _, f := range feedbacks {
		status := "Pending"
		if f.Resolved {
			status = "Resolved"
		}
		user := f.UserID
		if user == "" {
			user = "Anonymous"
		}
		fmt.Printf("ID:          %s\n", f.ID)
		fmt.Printf("Status:      %s\n", status)
		fmt.Printf("Type:        %s\n", f.Type)
		fmt.Printf("User:        %s\n", user)
		fmt.Printf("Fingerprint: %s\n", f.Fingerprint)
		fmt.Printf("Created:     %s\n", f.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Content:     %s\n", f.Content)
		fmt.Println("--------------------------------")
	}

	return nil
})
