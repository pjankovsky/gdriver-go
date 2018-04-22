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
	"crypto/subtle"
)

type FilePageData struct {
	CurrentPath  string
	Files        []*Path
	IsNotRoot    bool
	ParentPath   string
	ParentPathID string
}

func routers() *mux.Router {
	r := mux.NewRouter()

	// listing files - main display
	r.HandleFunc("/", handleFileList).Methods("GET")
	r.HandleFunc("/files/{path}", handleFileList).Methods("GET")

	// file queue handling
	r.HandleFunc("/queue/status", handleQueueStatus).Methods("POST")
	r.HandleFunc("/queue/update/{status}", handleQueueUpdate).Methods("POST")

	r.Use(authMiddleware)

	return r
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authOk := basicAuth(w, r)
		if authOk {
			next.ServeHTTP(w, r)
		}
	})
}

func basicAuth(w http.ResponseWriter, r *http.Request) bool {
	if settings.AuthUser == "" && settings.AuthPass == "" {
		return true
	}

	user, pass, ok := r.BasicAuth()

	if !ok ||
		subtle.ConstantTimeCompare([]byte(user), []byte(settings.AuthUser)) != 1 ||
		subtle.ConstantTimeCompare([]byte(pass), []byte(settings.AuthPass)) != 1 {
		w.Header().Set("WWW-Authenticate", `Basic realm="gdriver-go"`)
		w.WriteHeader(401)
		w.Write([]byte("Unauthorized.\n"))
		return false
	}
	return true
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

func parsePathPost(r *http.Request) ([]FileID, error) {
	decoder := json.NewDecoder(r.Body)
	var arr []FileID
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
		StatusReady,
		StatusPending,
		StatusInProgress,
		StatusDone:
		return status, nil
	}
	return "", fmt.Errorf("invalid argument: status")
}
