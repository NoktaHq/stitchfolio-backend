package requestModel

type FileStoreMetadata struct {
	ID       uint `json:"id,omitempty"`
	IsActive bool `json:"isActive,omitempty"`

	FileName   string `json:"fileName,omitempty"`
	FileSize   int64  `json:"fileSize,omitempty"`
	FileType   string `json:"fileType,omitempty"`
	FileUrl    string `json:"fileUrl,omitempty"`    // s3 url
	FileKey    string `json:"fileKey,omitempty"`    // s3 key
	FileBucket string `json:"fileBucket,omitempty"` // s3 bucket

	EntityId   uint   `json:"entityId,omitempty"`
	EntityType string `json:"entityType,omitempty"` // student, placement, fee, etc.
}
