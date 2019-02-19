package schema

// LocalFileID : DB autoinc ID
type LocalFileID uint64

// GDriveFileID : string ID from Google Drive
type GDriveFileID string

// MD5 : byte slice for MD5 hash
type MD5 [32]byte

// FileType : designate if the file is a Directory or File
type FileType string

// enum options for FileType
const (
	FileTypeFile FileType = "F"
	FileTypeDir  FileType = "D"
)

// UploadStatus : the status of the upload
type UploadStatus string

// enum options for UploadStatus
const (
	UploadStatusUnknown    UploadStatus = "unknown"
	UploadStatusError      UploadStatus = "error"
	UploadStatusReady      UploadStatus = "ready"
	UploadStatusPending    UploadStatus = "pending"
	UploadStatusInProgress UploadStatus = "inprogress"
	UploadStatusDone       UploadStatus = "done"
)
