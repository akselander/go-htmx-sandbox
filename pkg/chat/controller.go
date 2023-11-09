package chat

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Routes(g *echo.Group) {
	g.GET("/", index)
}

func index(c echo.Context) error {
	fmt.Printf("hey")
	return c.Render(http.StatusOK, "chat.html", "Hello")
}
