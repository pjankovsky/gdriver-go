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
)

type Path struct {
	ID      string
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
	StatusPending    Status = "pending"
	StatusInProgress Status = "inprogress"
	StatusDone       Status = "done"
)

func NewPath(file os.FileInfo, dirPath string) (*Path, error) {
	fulPath := path.Clean(dirPath + string(os.PathSeparator) + file.Name())
	relPath := fulPath[len(settings.Root):]
	fileID := base64.RawURLEncoding.EncodeToString([]byte(relPath))
	status, err := getFileStatus(fileID)
	if err != nil {
		return nil, err
	}

	p := &Path{
		ID:      fileID,
		Path:    fulPath,
		Name:    file.Name(),
		ModTime: file.ModTime(),
		IsDir:   file.IsDir(),
		Status:  status,
	}
	return p, nil
}

func listPaths(dirPath string) ([]*Path, error) {
	dirPath = path.Clean(dirPath)
	if strings.Contains(dirPath, settings.Root) == false {
		log.Printf("Target Path: %s", dirPath)
		return nil, errors.New("Path outside of defined root.")
	}

	fileInfoArr, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var pathArr []*Path
	for _, file := range fileInfoArr {
		path, err := NewPath(file, dirPath)
		if err != nil {
			return nil, err
		}
		pathArr = append(pathArr, path)
	}
	return pathArr, nil
}
