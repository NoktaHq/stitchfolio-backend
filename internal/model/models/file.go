package models

import "io"

type FileMetadata struct {
	Filename string
	Size     int64
	Header   map[string][]string
}

type FileUpload struct {
	Metadata FileMetadata  // Stores the FileName,Size and other metadata
	Content  io.ReadCloser // Stores the file content

	EntityId   uint
	EntityType string
	Kind       string
}

type ConfirmFile struct {
	Id          uint   `json:"id,omitempty"`
	FileKey     string `json:"fileKey,omitempty"`
	Kind        string `json:"kind,omitempty"`
	Description string `json:"description,omitempty"`
}

func (f *FileUpload) AddEntityInfo(entityId uint, entityType, kind string) bool {
	// If there is no content present , then return
	if f.Content == nil {
		return false
	}

	f.EntityId = entityId
	f.EntityType = entityType
	f.Kind = kind

	return true
}

func (f *FileUpload) HasContent() bool {
	return f != nil && f.Content != nil
}
