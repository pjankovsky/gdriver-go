package main

import (
	"log"
	"net/http"
	"sync"
)

func main() {
	loadSettings()
	getLocTime()
	setupBolt()
	getDrive() // checks tokens

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		r := routers()
		if settings.UseSsl == true {
			log.Fatal(http.ListenAndServeTLS(settings.IP+":"+settings.Port, sslCert, sslKey, r))
		} else {
			log.Fatal(http.ListenAndServe(settings.IP+":"+settings.Port, r))
		}
	}()

	go func() {
		defer wg.Done()
		waitAndUpload()
	}()

	wg.Wait()

}
