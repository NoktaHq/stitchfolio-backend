package responseModel

import "time"

type OrderHistory struct {
	ID                   uint       `json:"id,omitempty"`
	IsActive             bool       `json:"isActive,omitempty"`
	Action               string     `json:"action,omitempty"`
	ChangedFields        string     `json:"changedFields,omitempty"`
	Status               *string    `json:"status,omitempty"`
	ExpectedDeliveryDate *time.Time `json:"expectedDeliveryDate,omitempty"`
	DeliveredDate        *time.Time `json:"deliveredDate,omitempty"`
	OrderItemId          *uint      `json:"orderItemId,omitempty"`
	OrderItemData        string     `json:"orderItemData,omitempty"` // JSON string
	OrderId              uint       `json:"orderId,omitempty"`
	Order                *Order     `json:"order,omitempty"`
	PerformedAt          time.Time  `json:"performedAt,omitempty"`
	PerformedById        uint       `json:"performedById,omitempty"`
	PerformedBy          *User      `json:"performedBy,omitempty"`
}
