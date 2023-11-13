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
	"akselander/sandbox/pkg/shared/api"
	"akselander/sandbox/pkg/shared/database"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	port := flag.String("port", "", "Port application will run on")
	sessionSecret := flag.String("session", "", "Secret required to decypher sessions")
	databaseConnectionUri := flag.String("db", "", "Database connection string")
	flag.Parse()

	if *port == "" || *sessionSecret == "" || *databaseConnectionUri == "" {
		log.Fatalf("Could not initialize application, required flags missing.")
	}

	db, err := database.Connect(*databaseConnectionUri)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

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
	e.Use(api.ErrorHandlingMiddleware)
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(*sessionSecret))))
	e.Static("/dist", "dist")
	e.Static("/js", "js")

	e.GET("/", landing.Index)
	e.GET("/error", landing.Error)
	auth.App(db, e)
	chat.App(db, e)

	e.Logger.Fatal(e.Start(":" + *port))
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type TemplateRenderer struct {
	templates *template.Template
}
