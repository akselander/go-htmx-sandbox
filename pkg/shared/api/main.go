package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func ErrorHandlingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)

		if err == nil {
			return nil
		}

		template := EnsureProperTemplate(c, "error-card", "error-page")
		if c.Response().Status == http.StatusNotFound {
			c.Logger().Errorf("Not found: %v", c.Request().URL)
			template = EnsureProperTemplate(c, "not-found-card", "not-found-page")
		} else {
			c.Logger().Errorf("Internal error: %v", err)
		}

		return c.Render(http.StatusInternalServerError, template, nil)
	}
}

func IsHtmxRequest(c echo.Context) bool {
	r := c.Request().Header.Get("HX-Request")

	return r == "true"
}

func EnsureProperTemplate(c echo.Context, template string, fallback string) string {
	if IsHtmxRequest(c) {
		return template
	}

	return fallback
}
