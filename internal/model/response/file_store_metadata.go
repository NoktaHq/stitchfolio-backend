package responseModel

type FileStoreMetadata struct {
	Id       uint `json:"id,omitempty"`
	IsActive bool `json:"isActive,omitempty"`

	FileName   string `json:"fileName,omitempty"`
	FileSize   int64  `json:"fileSize,omitempty"`
	FileType   string `json:"fileType,omitempty"`
	FileUrl    string `json:"fileUrl,omitempty"`
	FileKey    string `json:"fileKey,omitempty"`
	FileBucket string `json:"fileBucket,omitempty"`

	EntityId   uint   `json:"entityId,omitempty"`
	EntityType string `json:"entityType,omitempty"`
}
