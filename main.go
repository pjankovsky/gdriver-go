package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

func main() {
	loadSettings()
	setupBolt()

	r := mux.NewRouter();

	// listing files - main display
	r.HandleFunc("/", handleFileList).Methods("GET")
	r.HandleFunc("/files/{path}", handleFileList).Methods("GET")

	// file queue handling
	r.HandleFunc("/queue/status", handleQueueStatus).Methods("POST")
	r.HandleFunc("/queue/update/{status}", handleQueueUpdate).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", r))
}
