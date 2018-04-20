package main

import (
	"net/http"
	"log"
	"fmt"
	"github.com/gorilla/mux"
	"encoding/json"
	"encoding/base64"
	"html/template"
	"path"
	"os"
)

type FilePageData struct {
	CurrentPath  string
	Files        []*Path
	IsNotRoot    bool
	ParentPath   string
	ParentPathID string
}

func handleFileList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	relPath := vars["path"]
	if len(relPath) > 0 {
		bytes, err := base64.RawURLEncoding.DecodeString(relPath)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), 500)
			return
		}
		relPath = string(bytes)
	} else {
		relPath = "/"
	}
	rootPath := path.Clean(settings.Root)

	fulPath := path.Clean(rootPath + string(os.PathSeparator) + relPath)

	paths, err := listPaths(fulPath)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}

	parentPath := path.Clean(relPath + string(os.PathSeparator) + "..")

	data := FilePageData{
		CurrentPath:  relPath,
		Files:        paths,
		IsNotRoot:    fulPath != rootPath,
		ParentPath:   parentPath,
		ParentPathID: base64.RawURLEncoding.EncodeToString([]byte(parentPath)),
	}

	tmpl, err := template.ParseFiles("html/files.html")
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}

	tmpl.Execute(w, data)
}

func handleQueueStatus(w http.ResponseWriter, r *http.Request) {
	fileIDs, err := parsePathPost(r)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}

	statusList, err := getFileStatusList(fileIDs)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(statusList)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}
}

func handleQueueUpdate(w http.ResponseWriter, r *http.Request) {
	fileIDs, err := parsePathPost(r)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}

	status, err := validateStatus(Status(mux.Vars(r)["status"]))
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}

	err = updateFileStatus(fileIDs, status)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}

	statusList, err := getFileStatusList(fileIDs)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(statusList)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}
}

func parsePathPost(r *http.Request) ([]string, error) {
	decoder := json.NewDecoder(r.Body)
	var arr []string
	err := decoder.Decode(&arr)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return arr, nil
}

func validateStatus(status Status) (Status, error) {
	switch status {
	case StatusUnknown,
		StatusError,
		StatusPending,
		StatusInProgress,
		StatusDone:
		return status, nil
	}
	return "", fmt.Errorf("Invalid argument: status")
}
