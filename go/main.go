package main

import (
	"html/template"
	"io"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	tmpl *template.Template
}

func newTemplate() *Template {
	return &Template{
		tmpl: template.Must(template.ParseGlob("views/*.html")),
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.tmpl.ExecuteTemplate(w, name, data)
}

type Contact struct {
	Name  string
	Email string
}

func newContact(name string, email string) Contact {
	return Contact{
		Name:  name,
		Email: email,
	}
}

type Count struct {
	Count int
}

type Result struct {
	PreviousValue string
	NewValue      string
	UnitFrom      string
	UnitTo        string
}

type Data struct {
	Count       Count
	Contacts    []Contact
	Form        FormData
	SelectedTab string
	Result      Result
}

func (d Data) IsTabSelected(name string) bool {
	return d.SelectedTab == name
}

func newData() Data {
	return Data{
		Contacts: []Contact{
			newContact("John Doe", "jd@email.com"),
			newContact("Claire Doe", "cd@email.com"),
		},
		Count:       Count{Count: 0},
		Form:        newFormData(),
		SelectedTab: "length",
		Result: Result{
			PreviousValue: "",
			NewValue:      "",
			UnitFrom:      "",
			UnitTo:        "",
		},
	}
}

type FormData struct {
	Values map[string]string
}

func newFormData() FormData {
	return FormData{Values: map[string]string{}}
}

func main() {

	e := echo.New()

	data := newData()

	e.Renderer = newTemplate()

	e.Use(middleware.Logger())

	e.Static("/static", "static")

	e.GET("/", func(c echo.Context) error {
		data.Count.Count++
		return c.Render(200, "index", data)
	})

	e.POST("/count", func(c echo.Context) error {
		data.Count.Count++
		return c.Render(200, "count", data)
	})

	e.POST("/contacts", func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")

		contact := newContact(name, email)
		data.Contacts = append(data.Contacts, contact)

		formData := newFormData()

		err := c.Render(200, "user-form", formData)

		if err != nil {
			return err
		}

		return c.Render(200, "oob-contact", contact)
	})

	e.GET("/form", func(c echo.Context) error {
		tab := c.QueryParam("tab")
		if tab == "" {
			tab = "length"
		}

		data.SelectedTab = tab

		return c.Render(200, "unit-converter", &data)

	})
	e.GET("/form/:tab", func(c echo.Context) error {
		tab := c.Param("tab")
		data.SelectedTab = tab
		err := c.Render(200, "tabs", &data)
		if err != nil {
			return err
		}

		return c.Render(200, "oob-form", &data)
	})

	e.POST("/form/:tab", func(c echo.Context) error {
		tab := c.Param("tab")
		data.Result.UnitFrom = c.FormValue("unit-from")
		data.Result.UnitTo = c.FormValue("unit-to")
		data.Result.PreviousValue = c.FormValue(tab)
		switch tab {
		case "length":
			length := c.FormValue("length")
			data.Result.NewValue = strconv.FormatFloat(computeLength(length, data.Result.UnitFrom, data.Result.UnitTo), 'f', -1, 64)
			return c.Render(200, "length-result", &data.Result)

		case "weight":
			weight := c.FormValue("weight")
			data.Result.NewValue = strconv.FormatFloat(computeWeight(weight, data.Result.UnitFrom, data.Result.UnitTo), 'f', -1, 64)
			return c.Render(200, "weight-result", &data.Result)

		case "temperature":
			temperature := c.FormValue("temperature")
			data.Result.NewValue = strconv.FormatFloat(computeTemperature(temperature, data.Result.UnitFrom, data.Result.UnitTo), 'f', -1, 64)
			return c.Render(200, "temperature-result", &data.Result)

		}

		return nil
	})

	e.Logger.Fatal(e.Start(":42069"))
}

func computeLength(length, unitFrom, unitTo string) float64 {
	val, err := strconv.ParseFloat(length, 64)
	if err != nil {
		return 0
	}

	switch unitFrom {
	case "cm":
		val = val / 100
	case "m":
		val = val * 1
	case "km":
		val = val * 1000
	}

	switch unitTo {
	case "cm":
		val = val * 100
	case "m":
		val = val * 1
	case "km":
		val = val / 1000
	}

	return val
}

func computeWeight(weight, unitFrom, unitTo string) float64 {
	val, err := strconv.ParseFloat(weight, 64)
	if err != nil {
		return 0
	}

	switch unitFrom {
	case "g":
		val = val / 1000
	case "kg":
		val = val * 1
	case "lb":
		val = val * 2.20462
	}

	switch unitTo {
	case "g":
		val = val * 1000
	case "kg":
		val = val * 1
	case "lb":
		val = val / 2.20462
	}
	return val
}

func computeTemperature(temperature, unitFrom, unitTo string) float64 {
	val, err := strconv.ParseFloat(temperature, 64)
	if err != nil {
		return 0
	}

	switch unitFrom {
	case "c":
		val = val * 1
	case "f":
		val = (val - 32) * 5 / 9
	}

	switch unitTo {
	case "c":
		val = val * 1
	case "f":
		val = val * 1
	}

	return val
}
