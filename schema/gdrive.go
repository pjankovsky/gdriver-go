package schema

import "time"

// GDriveFile is for db table gdrive_file
type GDriveFile struct {
	id          GDriveFileID
	parentID    GDriveFileID
	fileType    FileType
	name        string
	md5         MD5
	localFileID LocalFileID
	updatedAt   time.Time
	deleted     bool
}

// GDriveFileMetric is for db table gdrive_file_metric
type GDriveFileMetric struct {
	gDriveFileID GDriveFileID
	foundAt      time.Time
}
