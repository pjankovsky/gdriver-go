package main

import (
	"os"
	"io/ioutil"
	"time"
	"encoding/base64"
	"strings"
	"path"
	"errors"
	"log"
	"sort"
)

type FileID string

type Path struct {
	ID      FileID
	Path    string
	Name    string
	ModTime time.Time
	IsDir   bool
	Status  Status
}

type Status string

const (
	StatusUnknown    Status = "unknown"
	StatusError      Status = "error"
	StatusReady      Status = "ready"
	StatusPending    Status = "pending"
	StatusInProgress Status = "inprogress"
	StatusDone       Status = "done"
)

func NewPath(file os.FileInfo, dirPath string) (*Path, error) {
	fulPath := path.Clean(dirPath + string(os.PathSeparator) + file.Name())
	relPath := fulPath[len(settings.Root):]
	fileID := FileID(base64.RawURLEncoding.EncodeToString([]byte(relPath)))
	status, err := getFileStatus(fileID)
	if err != nil {
		return nil, err
	}

	p := &Path{
		ID:      fileID,
		Path:    fulPath,
		Name:    file.Name(),
		ModTime: file.ModTime().In(getLocTime()),
		IsDir:   file.IsDir(),
		Status:  status,
	}
	return p, nil
}

func getPath(relPath string) (*Path, error) {
	rootPath := path.Clean(settings.Root)
	fulPath := path.Clean(rootPath + string(os.PathSeparator) + relPath)

	fileInfo, err := os.Stat(fulPath)
	if err != nil {
		return nil, err
	}

	return NewPath(fileInfo, path.Dir(fulPath))
}

func listPaths(dirPath string) ([]*Path, error) {
	dirPath = path.Clean(dirPath)
	if strings.Contains(dirPath, settings.Root) == false {
		log.Printf("Target Path: %s", dirPath)
		return nil, errors.New("path outside of defined root")
	}

	fileInfoArr, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var pathArr []*Path
	for _, file := range fileInfoArr {
		pathSt, err := NewPath(file, dirPath)
		if err != nil {
			return nil, err
		}
		pathArr = append(pathArr, pathSt)
	}

	sort.Slice(pathArr, func(i, j int) bool {
		return pathArr[i].ModTime.UnixNano() > pathArr[j].ModTime.UnixNano()
	})

	return pathArr, nil
}

var locTime *time.Location

func getLocTime() *time.Location {
	if locTime != nil {
		return locTime
	}
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		log.Fatalf("Unable to local time location: %v", err)
	}
	locTime = loc
	return locTime
}
