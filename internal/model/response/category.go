package responseModel

type Category struct {
	ID       uint   `json:"id,omitempty"`
	IsActive bool   `json:"isActive,omitempty"`
	Name     string `json:"name,omitempty"`

	AuditFields `json:"auditFields,omitempty"`

	ProductCount int `json:"productCount,omitempty"` // Count of products in this category
}

type CategoryAutoComplete struct {
	ID   uint   `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
