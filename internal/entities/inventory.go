package entities

type Inventory struct {
	*Model `mapstructure:",squash"`

	ProductId         uint `json:"productId" gorm:"unique;not null"`
	Quantity          int  `json:"quantity" gorm:"not null;default:0"`
	LowStockThreshold int  `json:"lowStockThreshold" gorm:"default:0"`

	// Relations
	Product *Product `gorm:"foreignKey:ProductId" json:"product,omitempty"`
}

func (Inventory) TableNameForQuery() string {
	return "\"stich\".\"Inventories\" E"
}

// IsLowStock checks if current stock is below threshold
func (i *Inventory) IsLowStock() bool {
	return i.Quantity <= i.LowStockThreshold
}
