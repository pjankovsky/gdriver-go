package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"sync"
)

func main() {
	loadSettings()
	setupBolt()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		r := mux.NewRouter()

		// listing files - main display
		r.HandleFunc("/", handleFileList).Methods("GET")
		r.HandleFunc("/files/{path}", handleFileList).Methods("GET")

		// file queue handling
		r.HandleFunc("/queue/status", handleQueueStatus).Methods("POST")
		r.HandleFunc("/queue/update/{status}", handleQueueUpdate).Methods("POST")

		log.Print(http.ListenAndServe(":15445", r))
	}()

	go func() {
		defer wg.Done()
		waitAndUpload()
	}()

	wg.Wait()

}
