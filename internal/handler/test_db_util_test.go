package handler_test

import (
	"context"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/stretchr/testify/assert"
)

func TestSetupTestRepository(t *testing.T) {
	repo := handler.SetupTestRepository(t)
	assert.NotNil(t, repo)

	ctx := context.Background()

	// Verify we can save and read
	l := domain.Listing{
		ID:    "test-1",
		Title: "Test Listing",
		Type:  domain.Business,
	}

	err := repo.Save(ctx, l)
	assert.NoError(t, err)

	found, err := repo.FindByID(ctx, "test-1")
	assert.NoError(t, err)
	assert.Equal(t, "Test Listing", found.Title)
}
