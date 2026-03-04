package responseModel

type Product struct {
	ID           uint    `json:"id,omitempty"`
	IsActive     bool    `json:"isActive,omitempty"`
	Name         string  `json:"name,omitempty"`
	SKU          string  `json:"sku,omitempty"`
	CategoryId   uint    `json:"categoryId,omitempty"`
	Description  string  `json:"description,omitempty"`
	CostPrice    float64 `json:"costPrice,omitempty"`
	SellingPrice float64 `json:"sellingPrice,omitempty"`

	AuditFields `json:"auditFields,omitempty"`

	// Related data
	Category     *Category  `json:"category,omitempty"`
	Inventory    *Inventory `json:"inventory,omitempty"`
	CurrentStock int        `json:"currentStock,omitempty"` // From inventory
	IsLowStock   bool       `json:"isLowStock,omitempty"`   // Stock alert flag
	CategoryName string     `json:"categoryName,omitempty"` // Flattened category name
}

type ProductAutoComplete struct {
	ID           uint   `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	SKU          string `json:"sku,omitempty"`
	CurrentStock int    `json:"currentStock,omitempty"`
	IsLowStock   bool   `json:"isLowStock,omitempty"`
}
