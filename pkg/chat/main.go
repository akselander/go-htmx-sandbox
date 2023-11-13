package chat

import (
	"akselander/sandbox/pkg/auth"
	"database/sql"

	"net/http"

	"github.com/labstack/echo/v4"
)

func App(db *sql.DB, e *echo.Echo) {
	g := e.Group("/chat")
	g.Use(auth.LoginWallMiddleware)
	g.GET("/", index)
}

func index(c echo.Context) error {
	ac := c.(*auth.Context)
	return c.Render(http.StatusOK, "chat-page", Page{User: *ac.User})
}

type Page struct {
	User auth.User
}
