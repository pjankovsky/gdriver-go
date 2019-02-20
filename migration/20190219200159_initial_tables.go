package migration

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up20190219200159, Down20190219200159)
}

func Up20190219200159(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`
CREATE TABLE localFile (
  id INTEGER PRIMARY KEY,
  parentID INTEGER,
  fileType TEXT NOT NULL,
  name TEXT NOT NULL,
  md5 TEXT NOT NULL,
  gDriveFileID TEXT,
  uploadStatus TEXT NOT NULL DEFAULT "unknown",
  updatedAt TEXT NOT NULL,
  deleted INTEGER NOT NULL DEFAULT 0
);
`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
CREATE TABLE localFileMetric (
  localFileID INTEGER NOT NULL,
  foundAt TEXT,
  queuedAt TEXT,
  uploadStartedAt TEXT,
  uploadEndedAt TEXT,
  uploadRetries INTEGER NOT NULL DEFAULT 0
);
`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
CREATE TABLE gDriveFile (
  id TEXT PRIMARY KEY,
  parentID TEXT,
  fileType TEXT NOT NULL,
  name TEXT NOT NULL,
  md5 TEXT NOT NULL,
  localFileID INTEGER,
  updatedAt TEXT NOT NULL,
  deleted INTEGER NOT NULL DEFAULT 0
);
`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
CREATE TABLE gDriveFileMetric (
  gDriveFileID TEXT,
  foundAt TEXT
);
`)
	if err != nil {
		return err
	}
	return nil
}

func Down20190219200159(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	tx.Exec(`
DROP TABLE localFile;
DROP TABLE localFileMetric;
DROP TABLE gDriveFile;
DROP TABLE gDriveFileMetric;
`)
	return nil
}
