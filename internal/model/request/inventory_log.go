package requestModel

type InventoryLog struct {
	ID         uint   `json:"id,omitempty"`
	IsActive   bool   `json:"isActive,omitempty"`
	ProductId  uint   `json:"productId,omitempty"`
	ChangeType string `json:"changeType,omitempty"` // IN, OUT, ADJUST
	Quantity   int    `json:"quantity,omitempty"`
	Reason     string `json:"reason,omitempty"`
	Notes      string `json:"notes,omitempty"`
	LoggedAt   string `json:"loggedAt,omitempty"` // ISO datetime string
}

// StockMovementRequest is used for manual stock adjustments
type StockMovementRequest struct {
	ProductId     uint   `json:"productId" binding:"required"`
	ChangeType    string `json:"changeType" binding:"required"` // IN, OUT, ADJUST
	Quantity      int    `json:"quantity" binding:"required"`
	Reason        string `json:"reason" binding:"required"`
	Notes         string `json:"notes,omitempty"`
	AdminOverride bool   `json:"adminOverride,omitempty"` // Allow OUT even if stock insufficient
}
