package requestModel

import (
	"github.com/imkarthi24/sf-backend/internal/model/models"
)

type EntityDocuments struct {
	ID       uint `json:"id,omitempty"`
	IsActive bool `json:"isActive,omitempty"`

	Type         string `json:"type,omitempty"`
	DocumentType string `json:"documentType,omitempty"`
	Description  string `json:"description,omitempty"`
	EntityName   string `json:"entityName,omitempty"`
	EntityId     uint   `json:"entityId,omitempty"`

	Document *models.FileUpload `json:"-"`
}
