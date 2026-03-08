package responseModel

import "time"

type InventoryLog struct {
	ID         uint      `json:"id,omitempty"`
	IsActive   bool      `json:"isActive,omitempty"`
	ProductId  uint      `json:"productId,omitempty"`
	ChangeType string    `json:"changeType,omitempty"`
	Quantity   int       `json:"quantity,omitempty"`
	Reason     string    `json:"reason,omitempty"`
	Notes      string    `json:"notes,omitempty"`
	LoggedAt   time.Time `json:"loggedAt,omitempty"`
	SourceType string    `json:"sourceType,omitempty"` // MANUAL, PURCHASE, ADJUSTMENT
	SourceId   *uint     `json:"sourceId,omitempty"`

	AuditFields `json:"auditFields,omitempty"`

	// Related data
	Product      *Product `json:"product,omitempty"`
	ProductName  string   `json:"productName,omitempty"`
	ProductSKU   string   `json:"productSku,omitempty"`
	NetChange    string   `json:"netChange,omitempty"`    // e.g. "0(+20)" = 0 before, +20 added; "20(-5)" = 20 before, -5
	StockAfter   int      `json:"stockAfter,omitempty"`   // Stock quantity after this movement
	LoggedByName string   `json:"loggedByName,omitempty"` // User who logged
}

type StockMovementResponse struct {
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	ProductId     uint   `json:"productId"`
	PreviousStock int    `json:"previousStock"`
	NewStock      int    `json:"newStock"`
	ChangeAmount  int    `json:"changeAmount"`
}
