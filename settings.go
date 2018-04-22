package main

import (
	"os"
	"encoding/json"
	"path"
	"log"
)

type Settings struct {
	Root             string
	DbPath           string
	DriveRoot        string
	LocalTime        string
	UploadMaxWorkers int
	AuthUser         string
	AuthPass         string
}

var settings Settings

func loadSettings() {
	file, err := os.Open("etc/settings.json")
	if err != nil {
		log.Fatalf("Unable to load setting file: %v", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&settings)
	if err != nil {
		log.Fatalf("Unable to load setting file: %v", err)
	}
	settings.Root = path.Clean(settings.Root)

	log.Printf("Settings loaded:")
	log.Printf(" -- Root:              %v", settings.Root)
	log.Printf(" -- DbPath:            %v", settings.DbPath)
	log.Printf(" -- DriveRoot:         %v", settings.DriveRoot)
	log.Printf(" -- LocalTime:         %v", settings.LocalTime)
	log.Printf(" -- UploadMaxWorkers:  %v", settings.UploadMaxWorkers)
	log.Printf(" -- AuthUser:          %v", settings.AuthUser)
	log.Printf(" -- AuthPass:          %v", settings.AuthPass)
}
