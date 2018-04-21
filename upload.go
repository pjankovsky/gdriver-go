package main

import (
	"time"
	"encoding/base64"
	"log"
)

const (
	UploadTimeout    = 500 * time.Microsecond
	UploadMaxWorkers = 5
)

func waitAndUpload() {
	channel := make(chan FileID)
	defer close(channel)

	for i := 0; i < UploadMaxWorkers; i++ {
		go doUpload(channel)
	}

	for {
		fileID, err := claimFileForUpload()
		if err == nil {
			channel <- fileID
		}
		time.Sleep(UploadTimeout)
	}

}

func doUpload(in <-chan FileID) {
	for fileID := range in {
		if fileID == "" {
			continue
		}

		fileIDarr := make([]FileID, 1)
		fileIDarr[0] = fileID

		bytes, err := base64.RawURLEncoding.DecodeString(string(fileID))
		if err != nil {
			updateFileStatus(fileIDarr, StatusError)
			continue
		}

		relPath := string(bytes)
		path, err := getPath(relPath)
		if err != nil {
			updateFileStatus(fileIDarr, StatusError)
			continue
		}

		updateFileStatus(fileIDarr, StatusInProgress)

		log.Printf("file: %q", path.Name)
		time.Sleep(20 * time.Second)

		updateFileStatus(fileIDarr, StatusDone)
	}
}
