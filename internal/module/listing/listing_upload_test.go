package listing_test

import (
	listmod "github.com/jadecobra/agbalumo/internal/module/listing"

	"bytes"
	"context"
	"image"
	"image/png"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestListingHandler_Upload_Malicious(t *testing.T) {
	// Setup
	repo := handler.SetupTestRepository(t)
	listingSvc := listmod.NewListingService(repo, repo, repo)
	h := listmod.NewListingHandler(listmod.ListingDependencies{
		ListingStore:     repo,
		CategoryStore:    repo,
		ListingSvc:       listingSvc,
		ImageService:     nil,
		GeocodingSvc:     &MockGeocodingService{},
		Config:           &config.Config{},
	})

	// Create a malicious file (text file disguised as jpg)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Add File
	part, err := writer.CreateFormFile("image", "malicious.jpg")
	if err != nil {
		t.Fatal(err)
	}
	if _, err2 := part.Write([]byte("<?php echo 'malicious code'; ?>")); err2 != nil {
		t.Fatal(err2)
	}

	// Add Fields
	_ = writer.WriteField("title", "Valid Title")
	_ = writer.WriteField("owner_origin", "Nigeria")
	_ = writer.WriteField("type", "Business")
	_ = writer.WriteField("description", "Valid Description")
	_ = writer.WriteField("city", "Lagos")
	_ = writer.WriteField("address", "123 St")
	_ = writer.WriteField("contact_email", "test@test.com")

	_ = writer.Close()

	c, rec := setupTestContext(http.MethodPost, "/listings", body)
	c.Request().Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	c.Set("User", domain.User{ID: "user1", Email: "test@user.com"})

	// Execute
	_ = h.HandleCreate(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListingHandler_Upload_Valid(t *testing.T) {
	// Setup
	repo := handler.SetupTestRepository(t)
	listingSvc := listmod.NewListingService(repo, repo, repo)
	h := listmod.NewListingHandler(listmod.ListingDependencies{
		ListingStore:     repo,
		CategoryStore:    repo,
		ListingSvc:       listingSvc,
		ImageService:     nil,
		GeocodingSvc:     &MockGeocodingService{},
		Config:           &config.Config{},
	})

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Add valid PNG image
	part, _ := writer.CreateFormFile("image", "valid.png")
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	var imgBuf bytes.Buffer
	_ = png.Encode(&imgBuf, img)
	_, _ = part.Write(imgBuf.Bytes())

	// Add Fields
	_ = writer.WriteField("title", "Valid Title")
	_ = writer.WriteField("owner_origin", "Nigeria")
	_ = writer.WriteField("type", "Business")
	_ = writer.WriteField("description", "Valid Description")
	_ = writer.WriteField("city", "Lagos")
	_ = writer.WriteField("address", "123 St")
	_ = writer.WriteField("contact_email", "test@test.com")

	_ = writer.Close()

	c, rec := setupTestContext(http.MethodPost, "/listings", body)
	c.Request().Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	c.Set("User", domain.User{ID: "user1", Email: "test@user.com"})

	// Execute
	_ = h.HandleCreate(c)

	// Assert
	assert.Contains(t, []int{http.StatusOK, http.StatusCreated, http.StatusFound}, rec.Code)

	// Verify DB
	listings, _ := repo.FindByTitle(context.Background(), "Valid Title")
	assert.Len(t, listings, 1)
	assert.NotEmpty(t, listings[0].ImageURL)
}
