package mock_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/mock"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestMockListingRepository(t *testing.T) {
	m := &mock.MockListingRepository{}
	ctx := context.Background()

	// Test Save behavior
	m.On("Save", ctx, testifyMock.Anything).Return(nil).Once()
	if err := m.Save(ctx, domain.Listing{}); err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// Test Save Error
	m.On("Save", ctx, testifyMock.Anything).Return(errors.New("save error")).Once()
	if err := m.Save(ctx, domain.Listing{}); err == nil {
		t.Error("Expected error from Save")
	}

	// Test FindByID
	m.On("FindByID", ctx, "found").Return(domain.Listing{ID: "found"}, nil)
	if l, _ := m.FindByID(ctx, "found"); l.ID != "found" {
		t.Error("Expected ID to be 'found'")
	}

	// Test GetLocations
	m.On("GetLocations", ctx).Return([]string{"Lagos", "Accra"}, nil).Once()
	locations, err := m.GetLocations(ctx)
	if err != nil || len(locations) != 2 {
		t.Error("Expected 2 locations and nil error")
	}

	m.On("GetLocations", ctx).Return(nil, errors.New("location error")).Once()
	_, err = m.GetLocations(ctx)
	if err == nil {
		t.Error("Expected error from GetLocations")
	}

	// Assert
	m.AssertExpectations(t)
}

func TestMockListingRepository_AdditionalMethods(t *testing.T) {
	m := &mock.MockListingRepository{}
	ctx := context.Background()

	// FindAll
	m.On("FindAll", ctx, "type", "query", "", "", false, 10, 0).Return([]domain.Listing{{ID: "1"}}, nil).Once()
	listings, err := m.FindAll(ctx, "type", "query", "", "", false, 10, 0)
	if err != nil || len(listings) != 1 {
		t.Error("FindAll failed")
	}

	// FindByTitle
	m.On("FindByTitle", ctx, "title").Return([]domain.Listing{{ID: "1"}}, nil).Once()
	listings, err = m.FindByTitle(ctx, "title")
	if err != nil || len(listings) != 1 {
		t.Error("FindByTitle failed")
	}

	// SaveUser
	m.On("SaveUser", ctx, testifyMock.Anything).Return(nil).Once()
	err = m.SaveUser(ctx, domain.User{})
	if err != nil {
		t.Error("SaveUser failed")
	}

	// FindUserByGoogleID
	m.On("FindUserByGoogleID", ctx, "g1").Return(domain.User{ID: "u1"}, nil).Once()
	user, err := m.FindUserByGoogleID(ctx, "g1")
	if err != nil || user.ID != "u1" {
		t.Error("FindUserByGoogleID failed")
	}

	// FindUserByID
	m.On("FindUserByID", ctx, "u1").Return(domain.User{ID: "u1"}, nil).Once()
	user, err = m.FindUserByID(ctx, "u1")
	if err != nil || user.ID != "u1" {
		t.Error("FindUserByID failed")
	}

	// FindAllByOwner
	m.On("FindAllByOwner", ctx, "o1", 10, 0).Return([]domain.Listing{{ID: "1"}}, nil).Once()
	listings, err = m.FindAllByOwner(ctx, "o1", 10, 0)
	if err != nil || len(listings) != 1 {
		t.Error("FindAllByOwner failed")
	}

	// TitleExists
	m.On("TitleExists", ctx, "t1").Return(true, nil).Once()
	exists, err := m.TitleExists(ctx, "t1")
	if err != nil || !exists {
		t.Error("TitleExists failed")
	}

	// Delete
	m.On("Delete", ctx, "1").Return(nil).Once()
	err = m.Delete(ctx, "1")
	if err != nil {
		t.Error("Delete failed")
	}

	// GetCounts
	m.On("GetCounts", ctx).Return(map[domain.Category]int{domain.Business: 1}, nil).Once()
	counts, err := m.GetCounts(ctx)
	if err != nil || counts[domain.Business] != 1 {
		t.Error("GetCounts failed")
	}

	// ExpireListings
	m.On("ExpireListings", ctx).Return(int64(5), nil).Once()
	affected, err := m.ExpireListings(ctx)
	if err != nil || affected != 5 {
		t.Error("ExpireListings failed")
	}

	// SaveFeedback
	m.On("SaveFeedback", ctx, testifyMock.Anything).Return(nil).Once()
	err = m.SaveFeedback(ctx, domain.Feedback{})
	if err != nil {
		t.Error("SaveFeedback failed")
	}

	// GetAllFeedback
	m.On("GetAllFeedback", ctx).Return([]domain.Feedback{{ID: "1"}}, nil).Once()
	feedback, err := m.GetAllFeedback(ctx)
	if err != nil || len(feedback) != 1 {
		t.Error("GetAllFeedback failed")
	}

	// GetFeedbackCounts
	m.On("GetFeedbackCounts", ctx).Return(map[domain.FeedbackType]int{domain.FeedbackTypeIssue: 1}, nil).Once()
	fbCounts, err := m.GetFeedbackCounts(ctx)
	if err != nil || fbCounts[domain.FeedbackTypeIssue] != 1 {
		t.Error("GetFeedbackCounts failed")
	}

	// GetPendingClaimRequests
	m.On("GetPendingClaimRequests", ctx).Return([]domain.ClaimRequest{{ID: "cr1"}}, nil).Once()
	claims, err := m.GetPendingClaimRequests(ctx)
	if err != nil || len(claims) != 1 {
		t.Error("GetPendingClaimRequests failed")
	}

	// SaveClaimRequest
	m.On("SaveClaimRequest", ctx, testifyMock.Anything).Return(nil).Once()
	err = m.SaveClaimRequest(ctx, domain.ClaimRequest{})
	if err != nil {
		t.Error("SaveClaimRequest failed")
	}

	// UpdateClaimRequestStatus
	m.On("UpdateClaimRequestStatus", ctx, "cr1", domain.ClaimStatusApproved).Return(nil).Once()
	err = m.UpdateClaimRequestStatus(ctx, "cr1", domain.ClaimStatusApproved)
	if err != nil {
		t.Error("UpdateClaimRequestStatus failed")
	}

	// GetClaimRequestByUserAndListing
	m.On("GetClaimRequestByUserAndListing", ctx, "u1", "l1").Return(domain.ClaimRequest{ID: "cr1"}, nil).Once()
	claim, err := m.GetClaimRequestByUserAndListing(ctx, "u1", "l1")
	if err != nil || claim.ID != "cr1" {
		t.Error("GetClaimRequestByUserAndListing failed")
	}

	// GetUserCount
	m.On("GetUserCount", ctx).Return(100, nil).Once()
	count, err := m.GetUserCount(ctx)
	if err != nil || count != 100 {
		t.Error("GetUserCount failed")
	}

	// GetAllUsers
	m.On("GetAllUsers", ctx, 10, 0).Return([]domain.User{{ID: "1"}}, nil).Once()
	users, err := m.GetAllUsers(ctx, 10, 0)
	if err != nil || len(users) != 1 {
		t.Error("GetAllUsers failed")
	}

	// GetFeaturedListings
	m.On("GetFeaturedListings", ctx).Return([]domain.Listing{{ID: "1"}}, nil).Once()
	listings, err = m.GetFeaturedListings(ctx)
	if err != nil || len(listings) != 1 {
		t.Error("GetFeaturedListings failed")
	}

	// SetFeatured
	m.On("SetFeatured", ctx, "1", true).Return(nil).Once()
	err = m.SetFeatured(ctx, "1", true)
	if err != nil {
		t.Error("SetFeatured failed")
	}

	// GetListingGrowth
	m.On("GetListingGrowth", ctx).Return([]domain.DailyMetric{{Count: 1}}, nil).Once()
	growth, err := m.GetListingGrowth(ctx)
	if err != nil || len(growth) != 1 {
		t.Error("GetListingGrowth failed")
	}

	// GetUserGrowth
	m.On("GetUserGrowth", ctx).Return([]domain.DailyMetric{{Count: 1}}, nil).Once()
	userGrowth, err := m.GetUserGrowth(ctx)
	if err != nil || len(userGrowth) != 1 {
		t.Error("GetUserGrowth failed")
	}

	m.AssertExpectations(t)
}
