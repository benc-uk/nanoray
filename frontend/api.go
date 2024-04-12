package main

import (
	"fmt"
	"log"
	"nanoray/lib/controller"
	"nanoray/lib/proto"
	"net/http"
	"strconv"

	"google.golang.org/protobuf/types/known/wrapperspb"
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

	mux.HandleFunc("POST /api/render", func(w http.ResponseWriter, r *http.Request) {
		sceneData := r.FormValue("sceneData")

		width, _ := strconv.Atoi(r.FormValue("width"))
		aspectRatio, _ := strconv.ParseFloat(r.FormValue("aspect"), 64)
		height := int(float64(width) / aspectRatio)
		samplesPerPixel, _ := strconv.Atoi(r.FormValue("samples"))

		resp, err := controller.Client.StartRender(r.Context(), &proto.RenderRequest{
			SceneData: sceneData,
			ImageDetails: &proto.ImageDetails{
				Width:       int32(width),
				Height:      int32(height),
				AspectRatio: aspectRatio,
			},
			SamplesPerPixel: int32(samplesPerPixel),
			MaxDepth:        5,
			ChunkSize:       1,
		})

		if err != nil {
			log.Println("Failed to start render: ", err)
			http.Error(w, "Render failed to start: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte(fmt.Sprintf("Render started: %s", resp)))
	})

	mux.HandleFunc("GET /api/render", func(w http.ResponseWriter, r *http.Request) {
		data, _ := controller.Client.ListRenderedImages(r.Context(), nil)

		log.Println(data.Images)
		templates.Render(w, "api/renders", data)
	})

	mux.HandleFunc("GET /api/render/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		data, err := controller.Client.GetRenderedImage(r.Context(), wrapperspb.String(name))
		if err != nil {
			http.Error(w, "Failed to get rendered image "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.Write(data.Value)

	})
}
