package main

import (
	"html/template"
	"io"
	"strconv"
	"time"

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

var id = 0

type Contact struct {
	Name  string
	Email string
    Id int
}

func newContact(name string, email string) Contact {
    id++
	return Contact{
		Name:  name,
		Email: email,
        Id: id,
	}
}

type Contacts = []Contact

type Data struct {
	Contacts Contacts
}

func (d *Data) emailExists(email string) bool {
	for _, c := range d.Contacts {
		if c.Email == email {
			return true
		}
	}
	return false
}

func (d *Data) indexOf(id int) int {
    for i, contact := range d.Contacts {
        if contact.Id == id {
            return i
        }
    }

    return -1
}

func newData() Data {
	return Data{
		Contacts: []Contact{
			newContact("Steve", "steve@baboonmail.com"),
			newContact("Sam", "sam@proton.me"),
		},
	}
}

type FormData struct {
	Values map[string]string
	Errors map[string]string
}

func newFormData() FormData {
	return FormData{
		Values: make(map[string]string),
		Errors: make(map[string]string),
	}
}

type Page struct {
	Data Data
	Form FormData
}

func newPage() Page {
	return Page{
		Data: newData(),
		Form: newFormData(),
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.Renderer = newTempalte()

	page := newPage()

    e.Static("/images", "images")
    e.Static("/css", "css")

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", page)
	})

	e.POST("/contacts", func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")

		if page.Data.emailExists(email) {
			formData := newFormData()
			formData.Values["name"] = name
			formData.Values["email"] = email
			formData.Errors["email"] = "Email already exists"

			return c.Render(422, "form", formData)
		}

        contact := newContact(name, email)
		page.Data.Contacts = append(page.Data.Contacts, contact)

        err := c.Render(200, "form", newFormData())

        if err != nil {
            return err
        }
		return c.Render(200, "oob_contact", contact)
	})

    e.DELETE("/contacts/:id", func(c echo.Context) error {
        time.Sleep(3 * time.Second)
        idstr := c.Param("id")
        id, err := strconv.Atoi(idstr)

        if err != nil {
            return c.String(400, "Invalid id")
        }

        index := page.Data.indexOf(id)

        if index == -1 {
            return c.String(400, "Contact not found")
        }

        page.Data.Contacts = append(page.Data.Contacts[:index], page.Data.Contacts[index+1:]...)

        return c.NoContent(200)
    })

	e.Logger.Fatal(e.Start(":42069"))
}
