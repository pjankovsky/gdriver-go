package main

import (
	"os"
	"encoding/json"
	"path"
	"log"
)

type Settings struct {
	Root      string
	DbPath    string
	DriveRoot string
	AuthUser  string
	AuthPass  string
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
}
