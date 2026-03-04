package responseModel

type DressType struct {
	ID       uint `json:"id,omitempty"`
	IsActive bool `json:"isActive,omitempty"`

	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	Measurements string `json:"measurements,omitempty"`

	AuditFields `json:"auditFields,omitempty"`
}
