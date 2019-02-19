package schema

type LocalFileId int

type GDriveFileId string

type MD5 string

type FileType string

const (
	FileTypeFile FileType = "F"
	FileTypeDir  FileType = "D"
)

type Status string

const (
	StatusUnknown    Status = "unknown"
	StatusError      Status = "error"
	StatusReady      Status = "ready"
	StatusPending    Status = "pending"
	StatusInProgress Status = "inprogress"
	StatusDone       Status = "done"
)
