package handler_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/labstack/echo/v4"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestListingHandler_Upload_Malicious(t *testing.T) {
	// Setup
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	mockRepo := &mock.MockListingRepository{}
	// Save is called due to RespondError returning nil (likely), so execution continues.
	// We allow it here to prevent panic, but reliance is on rec.Code == 400.
	mockRepo.On("FindByTitle", testifyMock.Anything, testifyMock.Anything).Return([]domain.Listing{}, nil).Maybe()
	mockRepo.On("Save", testifyMock.Anything, testifyMock.Anything).Return(nil)

	h := handler.NewListingHandler(mockRepo, nil)

	// Create a malicious file (text file disguised as jpg)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Add File
	part, err := writer.CreateFormFile("image", "malicious.jpg")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write([]byte("<?php echo 'malicious code'; ?>")); err != nil {
		t.Fatal(err)
	}

	// Add Fields
	writer.WriteField("title", "Valid Title")
	writer.WriteField("owner_origin", "Nigeria")
	writer.WriteField("type", "Business")
	writer.WriteField("description", "Valid Description")
	writer.WriteField("city", "Lagos")
	writer.WriteField("address", "123 St")
	writer.WriteField("contact_email", "test@test.com")
	writer.WriteField("created_at", time.Now().Format(time.RFC3339))

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/listings", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("User", domain.User{ID: "user1", Email: "test@user.com"})

	// Execute
	err = h.HandleCreate(c)

	// Assert
	if rec.Code != http.StatusBadRequest {
		// If code is not 400, it means it succeeded (created).
		// If file was ignored, it created a listing without image (200 OK).
		// This confirms file was ignored.
		t.Errorf("Expected 400 Bad Request. Got %d. This implies file was ignored or accepted.", rec.Code)
	}
	mockRepo.AssertExpectations(t)
}

func TestListingHandler_Upload_Valid(t *testing.T) {
	// Setup
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	mockRepo := &mock.MockListingRepository{}
	// Save SHOULD be called
	mockRepo.On("FindByTitle", testifyMock.Anything, testifyMock.Anything).Return([]domain.Listing{}, nil).Maybe()
	mockRepo.On("Save", testifyMock.Anything, testifyMock.Anything).Return(nil)

	h := handler.NewListingHandler(mockRepo, nil)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Add Valid GIF/PNG/JPG magic bytes
	part, _ := writer.CreateFormFile("image", "valid.png")
	// tiny png signature (8 bytes)
	part.Write([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A})

	// Add Fields
	writer.WriteField("title", "Valid Title")
	writer.WriteField("owner_origin", "Nigeria")
	writer.WriteField("type", "Business")
	writer.WriteField("description", "Valid Description")
	writer.WriteField("city", "Lagos")
	writer.WriteField("address", "123 St")
	writer.WriteField("contact_email", "test@test.com")

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/listings", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("User", domain.User{ID: "user1", Email: "test@user.com"})

	// Execute
	h.HandleCreate(c)

	// Assert
	if rec.Code != http.StatusOK && rec.Code != http.StatusCreated && rec.Code != http.StatusFound {
		t.Errorf("Expected 200/201/302 for valid file, got %d. Body: %s", rec.Code, rec.Body.String())
	}
	mockRepo.AssertExpectations(t)
}
