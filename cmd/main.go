package main

import (
	"html/template"
	"io"
	"log"

	"akselander/sandbox/pkg/chat"
	"akselander/sandbox/pkg/landing"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	tmpls, err := template.New("").ParseGlob("pkg/*/views/*.html")

	if err != nil {
		log.Fatalf("Could not initialize templates: %v", err)
	}

	e := echo.New()
	e.Renderer = &TemplateRenderer{
		templates: tmpls,
	}

	e.Use(middleware.Logger())
	e.Static("/dist", "dist")
	e.Static("/css", "css")
	e.Static("/htmx", "htmx")

	e.GET("/", landing.Index)
	chat.Routes(e.Group("/chat"))

	e.Logger.Fatal(e.Start(":2137"))
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type TemplateRenderer struct {
	templates *template.Template
}
