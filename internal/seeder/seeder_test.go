package seeder_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/jadecobra/agbalumo/internal/seeder"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestSeedAll(t *testing.T) {
	mockRepo := &mock.MockListingRepository{}
	// Expect FindAll checks before seeding
	mockRepo.On("FindAll", context.Background(), testifyMock.Anything, testifyMock.Anything, testifyMock.Anything, 1, 0).Return([]domain.Listing{}, nil)
	// Expect Save calls - we just use Anything for simplicity as exact count varies
	mockRepo.On("Save", context.Background(), testifyMock.Anything).Return(nil)

	seeder.SeedAll(context.Background(), mockRepo)

	// Since we are mocking, we can't easily count actual logic calls without logic inside mock unless we setup specific expectations
	// For this test, verifying it runs without panic and calls Save is good enough proxy
	mockRepo.AssertCalled(t, "Save", context.Background(), testifyMock.Anything)
}

func TestEnsureSeeded_Empty(t *testing.T) {
	mockRepo := &mock.MockListingRepository{}
	// FindAll returns empty -> proceed to seed
	mockRepo.On("FindAll", context.Background(), testifyMock.Anything, testifyMock.Anything, testifyMock.Anything, 1, 0).Return([]domain.Listing{}, nil)
	mockRepo.On("Save", context.Background(), testifyMock.Anything).Return(nil)

	seeder.EnsureSeeded(context.Background(), mockRepo)

	mockRepo.AssertCalled(t, "Save", context.Background(), testifyMock.Anything)
}

func TestEnsureSeeded_NotEmpty(t *testing.T) {
	mockRepo := &mock.MockListingRepository{}
	// FindAll returns something -> skip seed
	mockRepo.On("FindAll", context.Background(), testifyMock.Anything, testifyMock.Anything, testifyMock.Anything, 1, 0).Return([]domain.Listing{{Title: "Existing"}}, nil)

	seeder.EnsureSeeded(context.Background(), mockRepo)

	mockRepo.AssertNotCalled(t, "Save", context.Background(), testifyMock.Anything)
}

func TestEnsureSeeded_FindAllError(t *testing.T) {
	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindAll", context.Background(), testifyMock.Anything, testifyMock.Anything, testifyMock.Anything, 1, 0).Return([]domain.Listing{}, errors.New("db error"))

	seeder.EnsureSeeded(context.Background(), mockRepo)

	mockRepo.AssertNotCalled(t, "Save", context.Background(), testifyMock.Anything)
}
