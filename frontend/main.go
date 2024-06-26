package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"nanoray/lib/controller"
)

var portFlag = flag.String("port", "8000", "The port to listen on")

func main() {
	portNum := os.Getenv("PORT")
	flag.Parse()

	if portNum == "" {
		portNum = *portFlag
	}
	port, _ := strconv.Atoi(portNum)

	fs := http.FileServer(http.Dir("public"))
	mux := http.NewServeMux()
	mux.Handle("/", fs)

	addAPIRoutes(mux, NewHTMLRenderer(mux))

	err := controller.Connect(time.Second * 20)
	if err != nil {
		log.Fatalf("Failed to connect to controller: %s", err.Error())
	}

	log.Printf("Frontend HTTP server started on port %d\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		log.Fatalf("Failed to start server\n%s", err.Error())
	}
}
