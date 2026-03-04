package requestModel

type Product struct {
	ID                uint    `json:"id,omitempty"`
	IsActive          bool    `json:"isActive,omitempty"`
	Name              string  `json:"name,omitempty"`
	SKU               string  `json:"sku,omitempty"`
	CategoryId        uint    `json:"categoryId,omitempty"`
	Description       string  `json:"description,omitempty"`
	CostPrice         float64 `json:"costPrice,omitempty"`
	SellingPrice      float64 `json:"sellingPrice,omitempty"`
	LowStockThreshold int     `json:"lowStockThreshold,omitempty"`
}
