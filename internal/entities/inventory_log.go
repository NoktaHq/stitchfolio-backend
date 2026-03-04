package entities

import "time"

type InventoryLogChangeType string

const (
	InventoryLogChangeTypeIN     InventoryLogChangeType = "IN"
	InventoryLogChangeTypeOUT    InventoryLogChangeType = "OUT"
	InventoryLogChangeTypeADJUST InventoryLogChangeType = "ADJUST"
)

type InventoryLog struct {
	*Model `mapstructure:",squash"`

	ProductId  uint                   `json:"productId" gorm:"not null"`
	ChangeType InventoryLogChangeType `json:"changeType" gorm:"type:varchar(20);not null"`
	Quantity   int                    `json:"quantity" gorm:"not null"`
	Reason     string                 `json:"reason" gorm:"not null"`
	Notes      string                 `json:"notes"`
	LoggedAt   time.Time              `json:"loggedAt" gorm:"not null"`

	// Relations
	Product *Product `gorm:"foreignKey:ProductId" json:"product,omitempty"`
}

func (InventoryLog) TableNameForQuery() string {
	return "\"stich\".\"InventoryLogs\" E"
}

// CalculateNetChange returns the net change in quantity based on change type
func (il *InventoryLog) CalculateNetChange() int {
	switch il.ChangeType {
	case InventoryLogChangeTypeIN:
		return il.Quantity
	case InventoryLogChangeTypeOUT:
		return -il.Quantity
	case InventoryLogChangeTypeADJUST:
		return il.Quantity // Can be positive or negative
	default:
		return 0
	}
}
