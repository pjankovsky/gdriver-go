package schema

import "time"

// GDriveFileRow is for db table gDriveFile
type GDriveFileRow struct {
	id          GDriveFileID
	parentID    GDriveFileID
	fileType    FileType
	name        string
	md5         MD5
	localFileID LocalFileID
	updatedAt   time.Time
	deleted     bool
}

// GDriveFileMetricRow is for db table gDriveFileMetric
type GDriveFileMetricRow struct {
	gDriveFileID GDriveFileID
	foundAt      time.Time
}
