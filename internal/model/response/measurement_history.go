package responseModel

import (
	"time"

	entitiy_types "github.com/imkarthi24/sf-backend/internal/entities/types"
)

type MeasurementHistory struct {
	ID            uint               `json:"id,omitempty"`
	IsActive      bool               `json:"isActive,omitempty"`
	Action        string             `json:"action,omitempty"`
	ChangedValues string             `json:"changedValues,omitempty"`
	OldValues     entitiy_types.JSON `json:"oldValues,omitempty"`
	MeasurementId uint               `json:"measurementId,omitempty"`
	Measurement   *Measurement       `json:"measurement,omitempty"`
	PerformedAt   time.Time          `json:"performedAt,omitempty"`
	PerformedById uint               `json:"performedById,omitempty"`
	PerformedBy   *User              `json:"performedBy,omitempty"`

	AuditFields `json:"auditFields,omitempty"`
}
