package requestModel

type Inventory struct {
	ID                uint `json:"id,omitempty"`
	IsActive          bool `json:"isActive,omitempty"`
	ProductId         uint `json:"productId,omitempty"`
	LowStockThreshold int  `json:"lowStockThreshold,omitempty"`
}
