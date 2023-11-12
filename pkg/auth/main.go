package auth

import (
	"database/sql"

	"net/http"

	"github.com/labstack/echo/v4"
)

func App(db *sql.DB, e *echo.Echo) {
	g := e.Group("/auth")
	g.GET("/", redirect)

	routes(g.Group("/signup"), "signup")
	routes(g.Group("/login"), "login")
}

func redirect(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/auth/signup/")
}

func routes(g *echo.Group, pageType string) {
	g.GET("/", page(pageType, "auth-page"))
	g.POST("/", page(pageType, "auth-username"))
	g.POST("/password/", page(pageType, "auth-password"))
}

func page(pageType string, template string) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Add("Cache-control", "no-cache, private")
		user := getUser(c)

		return c.Render(http.StatusOK, template, AuthPage{User: user, Page: pageType})
	}
}

func getUser(c echo.Context) User {
	username := c.FormValue("username")
	password := c.FormValue("password")

	return User{Username: username, Password: password}
}

type AuthPage struct {
	Page   string
	User   User
	Errors map[string]string
}

type User struct {
	Username string
	Password string
}
