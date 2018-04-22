package main

import (
	"time"
	"encoding/base64"
	"google.golang.org/api/drive/v3"
	"os"
	"log"
	"fmt"
)

const (
	UploadTimeout = 500 * time.Microsecond
)

func waitAndUpload() {
	channel := make(chan FileID)
	defer close(channel)

	for i := 0; i < settings.UploadMaxWorkers; i++ {
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
			log.Printf("Error decoding fileID [[%s]]: %v", fileID, err)
			updateFileStatus(fileIDarr, StatusError)
			continue
		}

		relPath := string(bytes)
		path, err := getPath(relPath)
		if err != nil {
			log.Printf("Error finding fileID [[%s]]: %v", fileID, err)
			updateFileStatus(fileIDarr, StatusError)
			continue
		}

		updateFileStatus(fileIDarr, StatusInProgress)

		err = uploadFileOrFolder(path, settings.DriveRoot, "")
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

	_, err = getDrive().Files.Create(&file).Media(b).Do()

	if err != nil {
		log.Printf("%s/%s]] ERROR: %v", parentLog, path.Name, err)
	}
	log.Printf("%s/%s]] DONE", parentLog, path.Name)

	return err
}
