package schema

import "time"

// LocalFileRow is for db table localFile
type LocalFileRow struct {
	id           LocalFileID
	parentID     LocalFileID
	fileType     FileType
	name         string
	md5          MD5
	gDriveFileID GDriveFileID
	uploadStatus UploadStatus
	updatedAt    time.Time
	deleted      bool
}

// LocalFileMetricRow is for db table localFileMetric
type LocalFileMetricRow struct {
	localFileID     LocalFileID
	foundAt         time.Time
	queuedAt        time.Time
	uploadStartedAt time.Time
	uploadEndedAt   time.Time
	uploadRetries   int
}
