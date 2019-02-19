package schema

import "time"

type LocalFile struct {
	id           LocalFileId
	parentId     LocalFileId
	fileType     FileType
	name         string
	md5          MD5
	gDriveFileId GDriveFileId
	uploadStatus Status
	updatedAt    time.Time
	deleted      bool
}

type LocalFileMetric struct {
	localFileId     LocalFileId
	foundAt         time.Time
	queuedAt        time.Time
	uploadStartedAt time.Time
	uploadEndedAt   time.Time
	uploadRetries   int
}
