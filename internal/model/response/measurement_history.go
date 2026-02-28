package responseModel

import "time"

type MeasurementHistory struct {
	ID            uint            `json:"id,omitempty"`
	IsActive      bool            `json:"isActive,omitempty"`
	Action        string          `json:"action,omitempty"`
	ChangedValues string          `json:"changedValues,omitempty"`
	OldValues     RawMessage      `json:"oldValues,omitempty"`
	MeasurementId uint            `json:"measurementId,omitempty"`
	Measurement   *Measurement    `json:"measurement,omitempty"`
	PerformedAt   time.Time       `json:"performedAt,omitempty"`
	PerformedById uint            `json:"performedById,omitempty"`
	PerformedBy   *User           `json:"performedBy,omitempty"`

	AuditFields `json:"auditFields,omitempty"`
}
