package main

import (
	"os"
	"encoding/json"
	"path"
)

type Settings struct {
	Root   string
	DbPath string
}

var settings Settings

func loadSettings() error {
	file, err := os.Open("settings.json")
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&settings)
	settings.Root = path.Clean(settings.Root)
	return err
}
