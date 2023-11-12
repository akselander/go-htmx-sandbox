package main

import (
	"flag"
	"html/template"
	"io"
	"log"
	"net/http"

	"akselander/sandbox/pkg/auth"
	"akselander/sandbox/pkg/chat"
	"akselander/sandbox/pkg/landing"
	database "akselander/sandbox/pkg/shared"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	var Port = flag.String("port", "", "Port application will run on")
	var DatabaseConnectionUri = flag.String("db", "", "Database connection string")
	flag.Parse()

	if *Port == "" || *DatabaseConnectionUri == "" {
		log.Fatalf("Could not initialize application, required flags missing.")
	}

	var db = database.Connect(*DatabaseConnectionUri)

	tmpls, err := template.New("").ParseGlob("pkg/*/views/*.html")

	if err != nil {
		log.Fatalf("Could not initialize templates: %v", err)
	}

	e := echo.New()
	e.Renderer = &TemplateRenderer{
		templates: tmpls,
	}

	e.Use(middleware.Logger())
	e.Use(middleware.AddTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	}))
	e.Static("/dist", "dist")
	e.Static("/js", "js")

	e.GET("/", landing.Index)
	auth.App(db, e)
	chat.App(db, e)

	e.Logger.Fatal(e.Start(":" + *Port))
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type TemplateRenderer struct {
	templates *template.Template
}
