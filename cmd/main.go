package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

// This is the interface for a renderer in echo.
// Render(io.Writer, string, interface{}, Context) error

func (t *Templates) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTempalte() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Count struct {
	Count int
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.Renderer = newTempalte()

	count := Count{Count: 0}

	e.GET("/", func(c echo.Context) error {
		count.Count += 1
		return c.Render(200, "index", count)
	})

    e.Logger.Fatal(e.Start(":42069"))
}
