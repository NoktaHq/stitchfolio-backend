package entities

type PurchaseItem struct {
	*Model `mapstructure:",squash"`
	
	QuantityOrdered    int     `json:"quantityOrdered" gorm:"not null"`
	QuantityReceived   int     `json:"quantityReceived" gorm:"not null;default:0"`
	UnitCost           float64 `json:"unitCost" gorm:"type:decimal(10,2);not null"`
	LineTotal          float64 `json:"lineTotal" gorm:"type:decimal(12,2);default:0"`

	// Relations
	PurchaseId         uint    `json:"purchaseId" gorm:"not null"`
	Purchase *Purchase `gorm:"foreignKey:PurchaseId" json:"purchase,omitempty"`
	ProductId          uint    `json:"productId" gorm:"not null"`
	Product  *Product  `gorm:"foreignKey:ProductId" json:"product,omitempty"`
}

func (PurchaseItem) TableNameForQuery() string {
	return "\"stich\".\"PurchaseItems\" E"
}

func (PurchaseItem) TableName() string {
	return TableNameWithSchema("PurchaseItems")
}
