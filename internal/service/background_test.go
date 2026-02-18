package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/mock"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestBackgroundService_ExpireListings(t *testing.T) {
	// Setup Mock
	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("ExpireListings", context.Background()).Return(int64(5), nil)

	service := NewBackgroundService(mockRepo)

	// Since expireListings is private but we are in package service, we can call it.
	service.expireListings(context.Background())

	mockRepo.AssertExpectations(t)
}

func TestBackgroundService_ExpireListings_Error(t *testing.T) {
	// Setup Mock to return error
	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("ExpireListings", context.Background()).Return(int64(0), errors.New("db error"))

	service := NewBackgroundService(mockRepo)

	// Should not panic, just log error
	service.expireListings(context.Background())

	mockRepo.AssertExpectations(t)
}

func TestBackgroundService_StartTicker_Cancels(t *testing.T) {
	// Setup Mock
	mockRepo := &mock.MockListingRepository{}
	// It should be called at least once
	mockRepo.On("ExpireListings", testifyMock.Anything).Return(int64(0), nil)

	service := NewBackgroundService(mockRepo)
	ctx, cancel := context.WithCancel(context.Background())

	// Run Ticker in goroutine
	done := make(chan bool)
	go func() {
		service.StartTicker(ctx)
		done <- true
	}()

	// Allow it to run for a tiny bit to trigger the initial "run once"
	time.Sleep(10 * time.Millisecond)

	// Cancel context to stop ticker
	cancel()

	// Wait for return
	select {
	case <-done:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("StartTicker did not exit after context cancellation")
	}

	mockRepo.AssertCalled(t, "ExpireListings", testifyMock.Anything)
}
