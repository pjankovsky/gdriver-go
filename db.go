package main

import (
	"errors"
	"fmt"
	bolt "github.com/etcd-io/bbolt"
	"log"
)

const (
	FileQueueBucketName = "FileQueue"
)

func setupBolt() {
	db, err := bolt.Open(settings.DbPath, 0600, nil)
	if err != nil {
		log.Fatalf("Unable to open db: %v", err)
	}
	defer closeBolt(db)

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(FileQueueBucketName))
		if err != nil {
			return fmt.Errorf("create bucket %s failed: %v", FileQueueBucketName, err)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Unable to setup default buckets: %v", err)
	}
}

func closeBolt(db *bolt.DB) {
	err := db.Close()
	if err != nil {
		log.Fatalf("Unable to close bolt: %v", err)
	}
}

func claimFileForUpload() (FileID, error) {
	db, err := bolt.Open(settings.DbPath, 0600, nil)
	if err != nil {
		return "", err
	}
	defer closeBolt(db)

	var claimedFileID FileID

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(FileQueueBucketName))
		c := b.Cursor()

		// scan to find a pending file
		var foundFileID FileID
		for fileID, status := c.First(); fileID != nil; fileID, status = c.Next() {
			status, err := validateStatus(Status(status))
			if err != nil {
				continue
			}
			if status == StatusReady {
				foundFileID = FileID(fileID)
				break
			}
		}

		// none found, return a blank
		if foundFileID == "" {
			return errors.New("no fileID found to claim")
		}

		// mark it as inprogress
		err := b.Put([]byte(foundFileID), []byte(StatusPending))
		if err != nil {
			return err
		}

		claimedFileID = foundFileID
		return nil
	})

	return claimedFileID, err
}

func getFileStatus(fileID FileID) (Status, error) {
	var status Status

	db, err := bolt.Open(settings.DbPath, 0600, nil)
	if err != nil {
		return StatusError, err
	}
	defer closeBolt(db)

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(FileQueueBucketName))
		v := string(b.Get([]byte(fileID)))
		if v == "" {
			status = StatusUnknown
		} else {
			status = Status(v)
		}
		return nil
	})

	if err != nil {
		return StatusError, err
	}

	status, err = validateStatus(status)
	if err != nil {
		status = StatusUnknown
	}

	return status, nil
}

func getFileStatusList(fileIDs []FileID) (map[FileID]Status, error) {
	statusList := make(map[FileID]Status)

	for _, fileID := range fileIDs {
		status, err := getFileStatus(fileID)
		if err != nil {
			return nil, err
		}
		statusList[fileID] = status
	}

	return statusList, nil
}

func updateFileStatus(fileIDs []FileID, status Status) error {
	db, err := bolt.Open(settings.DbPath, 0600, nil)
	if err != nil {
		return err
	}
	defer closeBolt(db)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(FileQueueBucketName))

		for _, fileID := range fileIDs {
			err := b.Put([]byte(fileID), []byte(status))
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
