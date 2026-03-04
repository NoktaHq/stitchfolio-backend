package requestModel

type Category struct {
	ID       uint   `json:"id,omitempty"`
	IsActive bool   `json:"isActive,omitempty"`
	Name     string `json:"name,omitempty"`
}
