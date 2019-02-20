package main

import (
	"flag"
	"log"
	"net/http"
	"sync"
)

func main() {

	flag.Parse()

	switch mode := flag.Arg(0) ; mode {
	case "daemonV1":
		daemonV1()
	case "migrateToV2":
		migrateToV2()
	case "scanGDrive":
		scanGDrive()
	default:
		flag.PrintDefaults()
	}

}

func scanGDrive() {
	loadSettings()
	getLocTime()
	setupSQL()
	getDrive()

}

func migrateToV2() {

}

func daemonV1() {
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
