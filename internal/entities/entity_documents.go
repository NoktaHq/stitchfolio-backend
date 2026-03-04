package entities

type EntityDocumentsType string

const (
	CHECKLIST            EntityDocumentsType = "checklist"
	CANDIDATE_PROFILEPIC EntityDocumentsType = "candidate-profile-pic"
)

type EntityDocuments struct {
	*Model `mapstructure:",squash"`

	Type         EntityDocumentsType `json:"type,omitempty"`         // Should be one of the EntityDocumentsType so that we can categorize
	DocumentType string              `json:"documentType,omitempty"` // Could be random string , kind of like filename / identifier

	EntityName EntityName `json:"entityName,omitempty"`
	EntityId   uint       `json:"entityId,omitempty"`

	Description string `json:"description,omitempty"`
}

func (EntityDocuments) TableName() string {
	return "EntityDocuments"
}

func (EntityDocuments) TableNameForQuery() string {
	return "\"EntityDocuments\" E"
}
