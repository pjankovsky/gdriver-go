package main

import (
	"github.com/coreos/bbolt"
	"fmt"
)

func setupBolt() error {
	db, err := bolt.Open(settings.DbPath, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("FileQueue"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func getFileStatus(fileID string) (Status, error) {
	var status Status

	db, err := bolt.Open(settings.DbPath, 0600, nil)
	if err != nil {
		return StatusError, err
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("FileQueue"))
		v := string(b.Get([]byte(fileID)))
		if v == "" {
			status = StatusUnknown
		} else {
			status = Status(v)
		}
		return nil
	})

	if err != nil {
		return StatusError, err;
	}

	status, err = validateStatus(status)
	if err != nil {
		status = StatusUnknown
	}

	return status, nil
}

func getFileStatusList(fileIDs []string) (map[string]Status, error) {
	statusList := make(map[string]Status)

	for _, fileID := range fileIDs {
		status, err := getFileStatus(fileID)
		if err != nil {
			return nil, err
		}
		statusList[fileID] = status
	}

	return statusList, nil
}

func updateFileStatus(fileIDs []string, status Status) error {
	db, err := bolt.Open(settings.DbPath, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("FileQueue"))

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
