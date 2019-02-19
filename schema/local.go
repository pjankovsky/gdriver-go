package schema

import "time"

// LocalFile is for db table local_file
type LocalFile struct {
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

// LocalFileMetric is for db table local_file_metric
type LocalFileMetric struct {
	localFileID     LocalFileID
	foundAt         time.Time
	queuedAt        time.Time
	uploadStartedAt time.Time
	uploadEndedAt   time.Time
	uploadRetries   int
}
