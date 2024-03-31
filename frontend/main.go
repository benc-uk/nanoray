package main

import (
	"log"
	"net/http"

	pb "nanoray/pkg/proto"

	"google.golang.org/grpc"
)

var ctrlClient pb.ControllerClient

func main() {
	fs := http.FileServer(http.Dir("public"))
	mux := http.NewServeMux()
	mux.Handle("/", fs)

	mux.HandleFunc("/api/workers", func(w http.ResponseWriter, r *http.Request) {
		workers, err := ctrlClient.GetWorkers(r.Context(), nil)
		if err != nil {
			log.Printf("Failed to get workers\n%s", err.Error())
			http.Error(w, "Failed to get workers", http.StatusInternalServerError)
			return
		}

		html := ""
		for _, worker := range workers.Workers {
			html += `<li>` + worker.Id + ` ` + worker.Address + `</li>`
		}

		w.Write([]byte(html))

	})

	conn, err := grpc.Dial("localhost:5000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to controller\n%s", err.Error())
	}
	ctrlClient = pb.NewControllerClient(conn)

	log.Println("Listening...")
	http.ListenAndServe(":8000", mux)
}
