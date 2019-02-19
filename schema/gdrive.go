package schema

import "time"

type GDriveFile struct {
	id          GDriveFileId
	parentId    GDriveFileId
	fileType    FileType
	name        string
	md5         MD5
	localFileId LocalFileId
	updatedAt   time.Time
	deleted     bool
}

type GDriveFileMetric struct {
	gDriveFileId GDriveFileId
	foundAt      time.Time
}
