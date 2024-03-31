package main

import (
	"html/template"
	"io"
)

// Implement echo.Renderer interface
type HTMLRenderer struct {
	templates *template.Template
}

func (r *HTMLRenderer) Render(w io.Writer, name string, data interface{}) error {
	return r.templates.ExecuteTemplate(w, name, data)
}

// // This function returns a new HTMLRenderer which loads templates from
// // the ./templates directory
// func newHTMLRenderer() *HTMLRenderer {
// 	return &HTMLRenderer{
// 		templates: template.Must(template.ParseGlob("./templates/**/*.html")),
// 	}
// }

// // Register the view handlers
// func AddViewHandlers(e *echo.Echo) {
// 	g := e.Group("/view")

// 	// Views are sub sections of the application, like home or todos
// 	// These are rendered under the navbar
// 	g.GET("/:viewName", func(c echo.Context) error {

// 		name := c.Param("viewName")

// 		err := c.Render(http.StatusOK, "view/"+name, nil)

// 		if err != nil {
// 			return c.String(http.StatusNotFound, "Not Found")
// 		}

// 		return nil
// 	})
// }
