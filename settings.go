package main

import (
	"os"
	"encoding/json"
	"path"
	"log"
)

type Settings struct {
	Host             string
	Port             string
	AuthUser         string
	AuthPass         string
	LocalRoot        string
	DriveRoot        string
	DbPath           string
	LocalTime        string
	UploadMaxWorkers int
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
	settings.LocalRoot = path.Clean(settings.LocalRoot)

	log.Printf("Settings loaded:")
	log.Printf(" -- Host:              %v", settings.Host)
	log.Printf(" -- Port:              %v", settings.Port)
	log.Printf(" -- AuthUser:          %v", settings.AuthUser)
	log.Printf(" -- AuthPass:          %v", settings.AuthPass)
	log.Printf(" -- LocalRoot:         %v", settings.LocalRoot)
	log.Printf(" -- DriveRoot:         %v", settings.DriveRoot)
	log.Printf(" -- DbPath:            %v", settings.DbPath)
	log.Printf(" -- LocalTime:         %v", settings.LocalTime)
	log.Printf(" -- UploadMaxWorkers:  %v", settings.UploadMaxWorkers)
}
