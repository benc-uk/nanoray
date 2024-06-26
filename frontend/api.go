package main

import (
	"log"
	"nanoray/lib/controller"
	"nanoray/lib/proto"
	"net/http"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func addAPIRoutes(mux *http.ServeMux, templates *HTMLRenderer) {
	mux.HandleFunc("GET /api/workers", func(w http.ResponseWriter, r *http.Request) {
		data, err := controller.Client.GetWorkers(r.Context(), nil)
		if err != nil {
			http.Error(w, "Failed to get workers", http.StatusInternalServerError)
			return
		}

		_ = templates.Render(w, "api/workers", data)
	})

	mux.HandleFunc("GET /api/render/progress", func(w http.ResponseWriter, r *http.Request) {
		data, _ := controller.Client.GetProgress(r.Context(), nil)

		// Special case for when the render is complete
		if data.CompletedJobs == data.TotalJobs {
			_ = templates.Render(w, "api/render-end", data)
			return
		}

		_ = templates.Render(w, "api/render-progress", data)
	})

	mux.HandleFunc("POST /api/render", func(w http.ResponseWriter, r *http.Request) {
		sceneData := r.FormValue("sceneData")

		width, _ := strconv.Atoi(r.FormValue("width"))
		depth, _ := strconv.Atoi(r.FormValue("depth"))
		slices, _ := strconv.Atoi(r.FormValue("slices"))
		aspectRatio, _ := strconv.ParseFloat(r.FormValue("aspect"), 64)
		samplesPerPixel, _ := strconv.Atoi(r.FormValue("samples"))

		_, err := controller.Client.StartRender(r.Context(), &proto.RenderRequest{
			SceneData:       sceneData,
			Width:           int32(width),
			AspectRatio:     aspectRatio,
			SamplesPerPixel: int32(samplesPerPixel),
			MaxDepth:        int32(depth),
			Slices:          int32(slices),
		})

		if err != nil {
			log.Println("Failed to start render: ", err)
			http.Error(w, "Render failed to start: "+err.Error(), http.StatusInternalServerError)
			return
		}

		_ = templates.Render(w, "api/render-start", nil)
	})

	mux.HandleFunc("GET /api/render", func(w http.ResponseWriter, r *http.Request) {
		data, _ := controller.Client.ListRenderedImages(r.Context(), nil)

		_ = templates.Render(w, "api/renders", data)
	})

	mux.HandleFunc("GET /api/render/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		// Set max message size to 64MB
		maxSizeOption := grpc.MaxCallRecvMsgSize(64 * 10e6)
		data, err := controller.Client.GetRenderedImage(r.Context(), wrapperspb.String(name), maxSizeOption)
		if err != nil {
			http.Error(w, "Failed to get rendered image "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write(data.Value)
	})
}
