package admin_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/jadecobra/agbalumo/internal/module/admin"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/jadecobra/agbalumo/internal/ui"
	"github.com/stretchr/testify/assert"
)

func TestAdminHandler_HandleToggleFeatured_BadRequest_MissingID(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := admin.NewAdminHandler(env.App)

	c, rec := testutil.SetupAdminContext(http.MethodPost, "/admin/listings//featured", nil)
	// Missing ID param
	c.SetParamNames("id")
	c.SetParamValues("")

	err := h.HandleToggleFeatured(c)
	assert.NoError(t, err) // Echo handlers return nil and specify code in c.JSON
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp ui.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v. Body tracking: %s", err, rec.Body.String())
	}

	assert.Equal(t, "Listing ID is required", errResp.Error)
	assert.Equal(t, http.StatusBadRequest, errResp.Code)
}
