package responseModel

import (
	entitiy_types "github.com/imkarthi24/sf-backend/internal/entities/types"
)

type Measurement struct {
	ID       uint `json:"id,omitempty"`
	IsActive bool `json:"isActive,omitempty"`

	Values entitiy_types.JSON `json:"values,omitempty"`

	PersonId   *uint   `json:"personId,omitempty"`
	Person     *Person `json:"person,omitempty"`
	PersonName string  `json:"personName,omitempty"`

	DressTypeId *uint      `json:"dressTypeId,omitempty"`
	DressType   *DressType `json:"dressType,omitempty"`

	TakenById *uint  `json:"takenById,omitempty"`
	TakenBy   string `json:"takenBy,omitempty"` // first_name + last_name

	AuditFields `json:"auditFields,omitempty"`
}

type MeasurementBrowse struct {
	ID         uint   `json:"id,omitempty"`
	IsActive   bool   `json:"isActive,omitempty"`
	PersonName string `json:"personName,omitempty"`
	PersonId   uint   `json:"personId"`
	CustomerId uint   `json:"customerId"`
	DressTypes string `json:"dressTypes"` // CSV of DressType names
	UpdatedAt  string `json:"updatedAt,omitempty"`
	UpdatedBy  string `json:"updatedBy,omitempty"`
}
