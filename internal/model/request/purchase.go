package requestModel

import "time"

type Purchase struct {
	ID                    uint            `json:"id,omitempty"`
	IsActive              bool            `json:"isActive,omitempty"`
	VendorId              uint            `json:"vendorId,omitempty"`
	PurchaseNumber        string          `json:"purchaseNumber,omitempty"`
	PurchaseDate          time.Time       `json:"purchaseDate,omitempty"`
	Status                string          `json:"status,omitempty"`
	ExpectedDeliveryDate  *time.Time      `json:"expectedDeliveryDate,omitempty"`
	Notes                 string          `json:"notes,omitempty"`
	TotalAmount           float64         `json:"totalAmount,omitempty"`
	PaidAmount            float64         `json:"paidAmount,omitempty"`
	PaidAt                *time.Time      `json:"paidAt,omitempty"`
	PaymentMethod         string          `json:"paymentMethod,omitempty"`
	PurchaseItems         []PurchaseItem `json:"purchaseItems,omitempty"`
}

type PurchaseItem struct {
	ID               uint    `json:"id,omitempty"`
	IsActive         bool    `json:"isActive,omitempty"`
	PurchaseId       uint    `json:"purchaseId,omitempty"`
	ProductId        uint    `json:"productId,omitempty"`
	QuantityOrdered  int     `json:"quantityOrdered,omitempty"`
	QuantityReceived int     `json:"quantityReceived,omitempty"`
	UnitCost         float64 `json:"unitCost,omitempty"`
	LineTotal        float64 `json:"lineTotal,omitempty"`
}
