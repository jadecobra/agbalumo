package admin

import (
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module"
)

// AdminLoginViewModel is the typed data for the admin login view.
type AdminLoginViewModel struct {
	Error string
	module.BaseViewData
}

// AdminModalChartsViewModel is the typed data for the charts modal.
type AdminModalChartsViewModel struct {
	ListingGrowth []domain.DailyMetric
	UserGrowth    []domain.DailyMetric
	module.BaseViewData
}

// AdminModalUsersViewModel is the typed data for the user management modal.
type AdminModalUsersViewModel struct {
	Users []domain.User
	module.BaseViewData
	UserCount int
}

// AdminModalCategoryViewModel is the typed data for the category management modal.
type AdminModalCategoryViewModel struct {
	Categories []domain.CategoryData
	module.BaseViewData
}

// AdminModalModerationViewModel is the typed data for the claim moderation modal.
type AdminModalModerationViewModel struct {
	ClaimRequests []domain.ClaimRequest
	module.BaseViewData
}

// AdminModalBulkViewModel is the typed data for the bulk actions modal.
type AdminModalBulkViewModel struct {
	module.BaseViewData
}
