package responseModel

type MasterConfig struct {
	Id       uint `json:"id,omitempty"`
	IsActive bool `json:"isActive,omitempty"`

	Name          string `json:"name,omitempty"`
	Type          string `json:"type,omitempty"`
	CurrentValue  string `json:"currentValue,omitempty"`
	PreviousValue string `json:"previousValue,omitempty"`
	DefaultValue  string `json:"defaultValue,omitempty"`
	UseDefault    bool   `json:"useDefault,omitempty"`
	Description   string `json:"description,omitempty"`
	Format        string `json:"format,omitempty"`

	AuditFields `json:"auditFields,omitempty"`
}
