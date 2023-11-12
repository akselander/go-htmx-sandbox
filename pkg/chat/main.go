package chat

import (
	"database/sql"

	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func App(db *sql.DB, e *echo.Echo) {
	g := e.Group("/chat")
	g.GET("/", index)
}

func index(c echo.Context) error {
	fmt.Printf("hey")
	return c.Render(http.StatusOK, "chat-page", "Hello")
}
