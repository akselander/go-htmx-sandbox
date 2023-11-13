package landing

import (
	"akselander/sandbox/pkg/auth"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Index(c echo.Context) error {
	ac := c.(*auth.Context)
	return c.Render(http.StatusOK, "index-page", Page{User: *ac.User})
}

func Error(c echo.Context) error {
	ac := c.(*auth.Context)
	return c.Render(http.StatusInternalServerError, "error-page", Page{User: *ac.User})
}

func NotFound(c echo.Context) error {
	ac := c.(*auth.Context)
	return c.Render(http.StatusInternalServerError, "not-found-page", Page{User: *ac.User})
}

type Page struct {
	User auth.User
}
