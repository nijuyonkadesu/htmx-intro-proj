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

type Count struct {
	Count int
}

type Contact struct {
	Name  string
	Email string
}

// Takes the data and renders the html, with context
func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newContact(name, email string) Contact {
    return Contact{
        Name:  name,
        Email: email,
    }
}

type Contacts = []Contact
type Data struct {
    Contacts Contacts
}

func (d *Data) hasEmail(email string) bool {
    for _, c := range d.Contacts {
        if c.Email == email {
            return true
        }
    }
    return false
}

func newData() Data {
    return Data{
        Contacts: Contacts{
            newContact("John Doe", "jj@mail.com"),
            newContact("Kiko", "kk@gmail.com"),
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

// Parse out the html
func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
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

    page := newPage()

	count := Count{Count: 0}
	e.Renderer = newTemplate()

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", page)
	})

	e.POST("/count", func(c echo.Context) error {
		count.Count++
		return c.Render(200, "count", count)
	})

	e.POST("/contacts", func(c echo.Context) error {
        name := c.FormValue("name")
        email := c.FormValue("email")

        if page.Data.hasEmail(email) {
            formData := newFormData()
            formData.Values["name"] = name
            formData.Values["email"] = email 
            formData.Errors["email"] = "Email already exists"
            return c.Render(422, "form", formData)
        }

        page.Data.Contacts = append(page.Data.Contacts, newContact(name, email))
		return c.Render(200, "display", page.Data)
	})

	e.Logger.Fatal(e.Start(":42069"))
}
