package responseModel

import "time"

type Inventory struct {
	ID                uint      `json:"id,omitempty"`
	IsActive          bool      `json:"isActive,omitempty"`
	ProductId         uint      `json:"productId,omitempty"`
	Quantity          int       `json:"quantity,omitempty"`
	LowStockThreshold int       `json:"lowStockThreshold,omitempty"`
	UpdatedAt         time.Time `json:"updatedAt,omitempty"`

	AuditFields `json:"auditFields,omitempty"`

	// Related data
	Product     *Product `json:"product,omitempty"`
	ProductName string   `json:"productName,omitempty"`
	ProductSKU  string   `json:"productSku,omitempty"`
	IsLowStock  bool     `json:"isLowStock,omitempty"`
}

type LowStockItem struct {
	ProductId         uint   `json:"productId"`
	ProductName       string `json:"productName"`
	ProductSKU        string `json:"productSku"`
	CurrentStock      int    `json:"currentStock"`
	LowStockThreshold int    `json:"lowStockThreshold"`
	CategoryName      string `json:"categoryName,omitempty"`
}
