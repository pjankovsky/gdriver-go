package main

import (
	"code.cloudfoundry.org/bytefmt"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

type FileID string

type Path struct {
	ID          FileID
	Path        string
	Name        string
	ModTime     time.Time
	Size        int64
	DisplaySize string
	IsDir       bool
	Status      Status
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
	relPath := fulPath[len(settings.LocalRoot):]
	fileID := FileID(base64.RawURLEncoding.EncodeToString([]byte(relPath)))
	status, err := getFileStatus(fileID)
	if err != nil {
		return nil, err
	}

	p := &Path{
		ID:          fileID,
		Path:        fulPath,
		Name:        file.Name(),
		ModTime:     file.ModTime().In(getLocTime()),
		Size:        file.Size(),
		DisplaySize: bytefmt.ByteSize(uint64(file.Size())),
		IsDir:       file.IsDir(),
		Status:      status,
	}
	return p, nil
}

func getPath(relPath string) (*Path, error) {
	rootPath := path.Clean(settings.LocalRoot)
	fulPath := path.Clean(rootPath + string(os.PathSeparator) + relPath)

	fileInfo, err := os.Stat(fulPath)
	if err != nil {
		return nil, err
	}

	return NewPath(fileInfo, path.Dir(fulPath))
}

func listPaths(dirPath string) ([]*Path, error) {
	dirPath = path.Clean(dirPath)
	if strings.Contains(dirPath, settings.LocalRoot) == false {
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

		iDate := makeSortKey(pathArr[i])
		jDate := makeSortKey(pathArr[j])

		if iDate == jDate {
			// dates are equal, sort Name ASC
			return pathArr[i].Name < pathArr[j].Name
		}
		// sort Date DESC
		return iDate > jDate
	})

	return pathArr, nil
}

func makeSortKey(path *Path) string {
	return path.ModTime.Format("20060102")
}


