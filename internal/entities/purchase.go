package entities

import "time"

type PurchaseStatus string

const (
	PurchaseStatusDRAFT              PurchaseStatus = "DRAFT"
	PurchaseStatusORDERED            PurchaseStatus = "ORDERED"
	PurchaseStatusPARTIALLY_RECEIVED  PurchaseStatus = "PARTIALLY_RECEIVED"
	PurchaseStatusRECEIVED           PurchaseStatus = "RECEIVED"
	PurchaseStatusCANCELLED          PurchaseStatus = "CANCELLED"
)

type Purchase struct {
	*Model `mapstructure:",squash"`

	
	PurchaseNumber          string       `json:"purchaseNumber"`
	PurchaseDate            time.Time    `json:"purchaseDate" gorm:"not null"`
	Status                  PurchaseStatus `json:"status" gorm:"type:varchar(30);not null;default:DRAFT"`
	ExpectedDeliveryDate   *time.Time   `json:"expectedDeliveryDate,omitempty"`
	Notes                   string       `json:"notes" gorm:"type:text"`
	TotalAmount             float64      `json:"totalAmount" gorm:"type:decimal(12,2);default:0"`

	// Simple payment tracking (phase 2)
	PaidAmount   float64    `json:"paidAmount" gorm:"type:decimal(12,2);default:0"`
	PaidAt       *time.Time `json:"paidAt,omitempty"`
	PaymentMethod string    `json:"paymentMethod"`

	// Relations
	VendorId                uint         `json:"vendorId" gorm:"not null"`
	Vendor        *Vendor        `gorm:"foreignKey:VendorId" json:"vendor,omitempty"`
	PurchaseItems []PurchaseItem `gorm:"foreignKey:PurchaseId" json:"purchaseItems,omitempty"`
}

func (Purchase) TableNameForQuery() string {
	return "\"stich\".\"Purchases\" E"
}
