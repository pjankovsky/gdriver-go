package main

import (
	"encoding/base64"
	"fmt"
	"google.golang.org/api/drive/v3"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

const (
	UploadTimeout = 1 * time.Second
	BackoffMax    = 1 * time.Hour // max time to wait between retries
	RetryMax      = 12            // max count of retries at BackoffMax
)

func waitAndUpload() {
	files := make(chan FileID)
	ticker := time.Tick(UploadTimeout)

	defer close(files)

	for i := 0; i < settings.UploadMaxWorkers; i++ {
		go doUpload(i, files)
	}

	for {
		<-ticker
		fileID, err := claimFileForUpload()
		if err == nil {
			files <- fileID
		}
	}

}

func doUpload(id int, in <-chan FileID) {
	for fileID := range in {
		if fileID == "" {
			continue
		}

		workerPrefix := fmt.Sprintf("[Worker %v] [[", id)

		fileIDarr := make([]FileID, 1)
		fileIDarr[0] = fileID

		bytes, err := base64.RawURLEncoding.DecodeString(string(fileID))
		if err != nil {
			log.Printf("%v%s]] Error decoding fileID: %v", workerPrefix, fileID, err)
			updateFileStatus(fileIDarr, StatusError)
			continue
		}

		relPath := string(bytes)
		path, err := getPath(relPath)
		if err != nil {
			log.Printf("%v%s]]Error finding fileID: %v", workerPrefix, fileID, err)
			updateFileStatus(fileIDarr, StatusError)
			continue
		}

		updateFileStatus(fileIDarr, StatusInProgress)

		err = uploadFileOrFolder(path, settings.DriveRoot, workerPrefix)
		if err != nil {
			updateFileStatus(fileIDarr, StatusError)
			continue
		}

		updateFileStatus(fileIDarr, StatusDone)
	}
}

func uploadFileOrFolder(path *Path, parentId string, parentLog string) error {

	if parentLog == "" {
		parentLog = "[["
	}

	if path.IsDir != true {
		return uploadFile(path, parentId, parentLog)
	}

	parentLog = fmt.Sprintf("%s/%s", parentLog, path.Name)

	parentIds := make([]string, 1)
	parentIds[0] = parentId

	// make the folder
	file := drive.File{
		Name:     path.Name,
		Parents:  parentIds,
		MimeType: "application/vnd.google-apps.folder",
	}

	r, err := getDrive().Files.Create(&file).Fields("id").Do()
	if err != nil {
		return err
	}

	subPaths, err := listPaths(path.Path)
	if err != nil {
		return err
	}

	for _, subPath := range subPaths {
		err = uploadFileOrFolder(subPath, r.Id, parentLog)
		if err != nil {
			return err
		}
	}

	return nil
}

func uploadFile(path *Path, parentId string, parentLog string) error {

	b, err := os.Open(path.Path)
	if err != nil {
		return err
	}

	parentIds := make([]string, 1)
	parentIds[0] = parentId

	file := drive.File{
		Name:    path.Name,
		Parents: parentIds,
	}

	log.Printf("%s/%s]] START", parentLog, path.Name)

	var success = false
	var breakOut = false
	var tries = 0
	var maxTries = 0
	var delay = 0 * time.Second

	for success == false && breakOut == false {
		if delay > 0 {
			log.Printf("%s/%s]] TOTAL TRIES: %v", parentLog, path.Name, tries)
			log.Printf("%s/%s]] RETRY TIMEOUT: %s", parentLog, path.Name, delay)
		}
		time.Sleep(delay)

		_, err = getDrive().Files.Create(&file).Media(b).Do()

		if err != nil {
			log.Printf("%s/%s]] ERROR: %v", parentLog, path.Name, err)

			delay = calcBackoff(tries)
			if delay == BackoffMax {
				maxTries++
				if maxTries > RetryMax {
					breakOut = true
					log.Printf("%s/%s]] MAX RETRIES REACHED", parentLog, path.Name)
				}
			}
			tries++

		} else {
			success = true
		}
	}

	log.Printf("%s/%s]] DONE", parentLog, path.Name)
	return err
}

func calcBackoff(n int) time.Duration {
	backoff := time.Duration(math.Pow(2, float64(n))) * time.Second
	backoff += time.Duration(rand.Int31n(1000)) * time.Millisecond

	if backoff > BackoffMax {
		return BackoffMax
	}
	return backoff
}
