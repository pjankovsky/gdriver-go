package main

import (
	"log"
	"net/http"
	"sync"
)

func main() {
	getLocTime()
	loadSettings()
	setupBolt()
	getDrive() // checks tokens

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		r := routers()
		log.Fatal(http.ListenAndServe(":15445", r))
	}()

	go func() {
		defer wg.Done()
		waitAndUpload()
	}()

	wg.Wait()

}
