package mock

import (
	"io"

	"github.com/labstack/echo/v4"
)

type MockRenderer struct{}

func (m *MockRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return nil
}
