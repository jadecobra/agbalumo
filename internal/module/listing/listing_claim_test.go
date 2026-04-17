package listing_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module/listing"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestHandleClaim(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := listing.NewListingHandler(env.App)
	testutil.SaveTestListing(t, env.App.DB, "1", "Biz")
	_ = env.App.DB.SaveCategory(context.Background(), domain.CategoryData{ID: string(domain.Business), Name: string(domain.Business), Claimable: true, Active: true})

	c, rec := testutil.SetupModuleContext(http.MethodPost, "/listings/1/claim", nil)
	c.SetPath("/listings/:id/claim")
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("User", domain.User{ID: "claimer", Role: domain.UserRoleUser})

	_ = h.HandleClaim(c)

	assert.Equal(t, http.StatusOK, rec.Code)
}
