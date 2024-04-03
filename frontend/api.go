package main

import (
	"fmt"
	"nanoray/shared/controller"
	"net/http"
)

func addAPIRoutes(mux *http.ServeMux, templates *HTMLRenderer) {
	mux.HandleFunc("GET /api/workers", func(w http.ResponseWriter, r *http.Request) {
		data, err := controller.Client.GetWorkers(r.Context(), nil)
		if err != nil {
			http.Error(w, "Failed to get workers", http.StatusInternalServerError)
			return
		}

		templates.Render(w, "api/workers", data)
	})

	mux.HandleFunc("GET /api/render/progress", func(w http.ResponseWriter, r *http.Request) {
		data, _ := controller.Client.GetProgress(r.Context(), nil)

		templates.Render(w, "api/render-progress", data)
	})

	mux.HandleFunc("GET /api/render", func(w http.ResponseWriter, r *http.Request) {
		resp, err := controller.Client.StartRender(r.Context(), nil)
		if err != nil {
			http.Error(w, "Render failed to start: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte(fmt.Sprintf("Render started: %s", resp)))
	})
}
