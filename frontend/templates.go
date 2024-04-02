package main

import (
	"html/template"
	"io"
	"net/http"
)

// Implement echo.Renderer interface
type HTMLRenderer struct {
	templates *template.Template
}

func (r *HTMLRenderer) Render(w io.Writer, name string, data interface{}) error {
	return r.templates.ExecuteTemplate(w, name, data)
}

func NewHTMLRenderer(mux *http.ServeMux) *HTMLRenderer {
	rend := &HTMLRenderer{
		templates: template.Must(template.ParseGlob("templates/**/*.html")),
	}

	mux.HandleFunc("GET /view/{view}", func(w http.ResponseWriter, r *http.Request) {
		view := r.PathValue("view")
		err := rend.Render(w, "view/"+view, nil)
		if err != nil {
			http.Error(w, "Failed to render view", http.StatusInternalServerError)
		}
	})

	return rend
}
