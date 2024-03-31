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
