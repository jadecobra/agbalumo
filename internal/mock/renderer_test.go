package mock_test

import (
    "bytes"
	"testing"

	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestMockRenderer_Render(t *testing.T) {
	renderer := new(mock.MockRenderer)

    var b bytes.Buffer
    e := echo.New()
    c := e.NewContext(nil, nil)
    
	err := renderer.Render(&b, "test.html", nil, c)
	assert.NoError(t, err)
}
