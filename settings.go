package main

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"time"
)

type Settings struct {
	IP               string
	Port             string
	UseSsl           bool
	SslHostname      string
	AuthUser         string
	AuthPass         string
	LocalRoot        string
	DriveRoot        string
	DbPath           string
	LocalTime        string
	UploadMaxWorkers int
	SQLitePath       string
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

	setupSsl()

	log.Printf("Settings loaded:")
	log.Printf(" -- IP:                %v", settings.IP)
	log.Printf(" -- Port:              %v", settings.Port)
	log.Printf(" -- UseSsl:            %v", settings.UseSsl)
	log.Printf(" -- SslHostname:       %v", settings.SslHostname)
	log.Printf(" -- AuthUser:          %v", settings.AuthUser)
	log.Printf(" -- AuthPass:          %v", settings.AuthPass)
	log.Printf(" -- LocalRoot:         %v", settings.LocalRoot)
	log.Printf(" -- DriveRoot:         %v", settings.DriveRoot)
	log.Printf(" -- DbPath:            %v", settings.DbPath)
	log.Printf(" -- SQLitePath:        %v", settings.SQLitePath)
	log.Printf(" -- LocalTime:         %v", settings.LocalTime)
	log.Printf(" -- UploadMaxWorkers:  %v", settings.UploadMaxWorkers)
}

var locTime *time.Location

func getLocTime() *time.Location {
	if locTime != nil {
		return locTime
	}
	loc, err := time.LoadLocation(settings.LocalTime)
	if err != nil {
		log.Fatalf("Unable to local time location: %v", err)
	}
	locTime = loc
	return locTime
}
