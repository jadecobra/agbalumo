package handler_test

import (
	"bytes"
	"image"
	"image/png"
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

	h := handler.NewListingHandler(mockRepo, nil, "")

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
	_ = writer.WriteField("created_at", time.Now().Format(time.RFC3339))

	_ = writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/listings", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("User", domain.User{ID: "user1", Email: "test@user.com"})

	// Execute
	_ = h.HandleCreate(c)

	// Assert
	if rec.Code != http.StatusBadRequest {
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

	h := handler.NewListingHandler(mockRepo, nil, "")

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

	req := httptest.NewRequest(http.MethodPost, "/listings", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("User", domain.User{ID: "user1", Email: "test@user.com"})

	// Execute
	_ = h.HandleCreate(c)

	// Assert
	if rec.Code != http.StatusOK && rec.Code != http.StatusCreated && rec.Code != http.StatusFound {
		t.Errorf("Expected 200/201/302 for valid file, got %d. Body: %s", rec.Code, rec.Body.String())
	}
	mockRepo.AssertExpectations(t)
}
