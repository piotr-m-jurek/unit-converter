package main

import (
	"html/template"
	"io"
	"strconv"
	"unit-converter/domain"

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

type Result struct {
	PreviousValue string
	NewValue      string
	UnitFrom      string
	UnitTo        string
	Error         string
}

type Data struct {
	SelectedTab string
	Result      Result
}

func (d Data) IsTabSelected(name string) bool {
	return d.SelectedTab == name
}

func newData() Data {
	return Data{
		SelectedTab: "length",
		Result: Result{
			PreviousValue: "",
			NewValue:      "",
			UnitFrom:      "",
			UnitTo:        "",
			Error:         "",
		},
	}
}

func main() {

	e := echo.New()
	e.Renderer = newTemplate()
	e.Use(middleware.Logger())
	e.Static("/static", "static")

	data := newData()

	e.GET("/", func(c echo.Context) error {
		tab := c.QueryParam("tab")
		if tab == "" {
			tab = "length" // default tab
		}
		data.SelectedTab = tab

		if c.Request().Header.Get("HX-Request") == "true" {
			// If it's an HTMX request, just return the tabs and form
			err := c.Render(200, "tabs", &data)
			if err != nil {
				return err
			}
			return c.Render(200, "oob-form", &data)
		}

		// For regular requests, return the full page
		return c.Render(200, "unit-converter", &data)
	})

	e.POST("/", func(c echo.Context) error {
		tab := c.QueryParam("tab")
		data.Result.UnitFrom = c.FormValue("unit-from")
		data.Result.UnitTo = c.FormValue("unit-to")
		data.Result.PreviousValue = c.FormValue(tab)

		switch tab {
		case "length":
			length := c.FormValue("length")
			converter := domain.LengthConverter{}
			newValue, err := converter.Compute(length, data.Result.UnitFrom, data.Result.UnitTo)
			if err != nil {
				data.Result.Error = err.Error()
			}
			data.Result.NewValue = strconv.FormatFloat(newValue, 'f', -1, 64)
			return c.Render(200, "length-result", &data.Result)

		case "weight":
			weight := c.FormValue("weight")
			converter := domain.WeightConverter{}
			newValue, err := converter.Compute(weight, data.Result.UnitFrom, data.Result.UnitTo)
			if err != nil {
				data.Result.Error = err.Error()
			}
			data.Result.NewValue = strconv.FormatFloat(newValue, 'f', -1, 64)
			return c.Render(200, "weight-result", &data.Result)

		case "temperature":
			temperature := c.FormValue("temperature")
			converter := domain.TemperatureConverter{}
			newValue, err := converter.Compute(temperature, data.Result.UnitFrom, data.Result.UnitTo)
			if err != nil {
				data.Result.Error = err.Error()
			}
			data.Result.NewValue = strconv.FormatFloat(newValue, 'f', -1, 64)
			return c.Render(200, "temperature-result", &data.Result)

		}

		return nil
	})

	e.Logger.Fatal(e.Start(":42069"))
}
