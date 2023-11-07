package pages

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", "hellow")
}

func Todos(c echo.Context) error {
	return c.Render(http.StatusOK, "todos.html", "hellow")
}
