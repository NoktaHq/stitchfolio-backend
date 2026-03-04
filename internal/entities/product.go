package entities

type Product struct {
	*Model `mapstructure:",squash"`

	Name         string  `json:"name" gorm:"not null"`
	SKU          string  `json:"sku" gorm:"unique"`
	CategoryId   uint    `json:"categoryId" gorm:"not null"`
	Description  string  `json:"description" gorm:"type:text"`
	CostPrice    float64 `json:"costPrice" gorm:"type:decimal(10,2);not null"`
	SellingPrice float64 `json:"sellingPrice" gorm:"type:decimal(10,2);not null"`

	// Relations
	Category  *Category  `gorm:"foreignKey:CategoryId" json:"category,omitempty"`
	Inventory *Inventory `gorm:"foreignKey:ProductId" json:"inventory,omitempty"`
}

func (Product) TableNameForQuery() string {
	return "\"stich\".\"Products\" E"
}
