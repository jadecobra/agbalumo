package listing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestAdaDiscoveryFlow_Semantic(t *testing.T) {
	// 1. Setup Environment
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()

	// 2. Seed Data for Ada's Dallas Flow
	testutil.SeedAdaDallasData(t, env.App.DB)

	h := NewListingHandler(env.App)

	t.Run("Ada searches for Nigerian Food in Dallas", func(t *testing.T) {
		// Mock Request: GET /listings/fragment?type=Food&city=Dallas
		target := fmt.Sprintf("/listings/fragment?%s=%s&%s=%s",
			domain.FieldType, domain.Food,
			domain.FieldCity, "Dallas")

		c, rec := testutil.SetupModuleContext(http.MethodGet, target, nil)

		// 3. Execute Handler
		err := h.HandleFragment(c)
		assert.NoError(t, err)

		// 4. Semantic Assertions
		body := rec.Body.String()
		assert.Equal(t, http.StatusOK, rec.Code)

		// Verify Ada's specific listings are present using ag-test-id
		testutil.AssertContainsSemanticID(t, body, "listing-ada-1")
		testutil.AssertContainsSemanticID(t, body, "listing-ada-2")
		testutil.AssertContainsSemanticID(t, body, "listing-ada-3")

		// Verify Featured listing is correctly tagged
		assert.Contains(t, body, "Jollof House")
		assert.Contains(t, body, "Featured") // Simplified check for featured badge

		// Verify spatial context in result
		assert.Contains(t, body, "Dallas")
	})

	t.Run("Ada views detail of Nigerian restaurant", func(t *testing.T) {
		// Mock Request: GET /listings/ada-1
		target := "/listings/ada-1"
		c, rec := testutil.SetupModuleContext(http.MethodGet, target, nil)
		c.SetParamNames("id")
		c.SetParamValues("ada-1")

		// Execute Handler
		err := h.HandleDetail(c)
		assert.NoError(t, err)

		// Semantic Assertions
		body := rec.Body.String()
		assert.Equal(t, http.StatusOK, rec.Code)

		// Verify semantic tags in modal
		testutil.AssertContainsSemanticID(t, body, "modal-detail")
		assert.Contains(t, body, "Jollof House")
		assert.Contains(t, body, "Authentic Nigerian Jollof")

		// Verify data agent template metadata (if present)
		assert.Contains(t, body, "data-agent-template=\"modal_detail\"")
	})
}
