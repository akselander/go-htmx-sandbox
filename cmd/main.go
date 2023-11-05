package main

import (
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	tmpls, err := template.New("").ParseGlob("public/views/*.html")

	if err != nil {
		log.Fatalf("Could not initialize templates: %v", err)
	}

	e := echo.New()
	e.Renderer = &TemplateRenderer{
		templates: tmpls,
	}

	e.Use(middleware.Logger())
	e.Static("/htmx", "htmx")

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", "Sandbox time!")
	})

	e.Logger.Fatal(e.Start(":42069"))
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type TemplateRenderer struct {
	templates *template.Template
}
