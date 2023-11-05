package main

import (
	"context"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", root)
	e.Logger.Fatal(e.Start(":42069"))
}

func root(c echo.Context) error {
	return view().Render(context.Background(), c.Response().Writer)
}
