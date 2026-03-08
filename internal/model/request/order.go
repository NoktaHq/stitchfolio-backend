package requestModel

import "github.com/imkarthi24/sf-backend/internal/model/models"

type Order struct {
	ID       uint `json:"id,omitempty"`
	IsActive bool `json:"isActive,omitempty"`

	Status string `json:"status,omitempty"`

	Notes string `json:"notes,omitempty"`

	AdditionalCharges float64 `json:"additionalCharges,omitempty"`

	ExpectedDeliveryDate *string `json:"expectedDeliveryDate,omitempty"`
	DeliveredDate        *string `json:"deliveredDate,omitempty"`

	CustomerId     *uint `json:"customerId,omitempty"`
	OrderTakenById *uint `json:"orderTakenById,omitempty"`

	OrderItems []OrderItem `json:"orderItems,omitempty"`
}

type OrderItem struct {
	ID       uint `json:"id,omitempty"`
	IsActive bool `json:"isActive,omitempty"`

	Description       string  `json:"description,omitempty"`
	Quantity          int     `json:"quantity,omitempty"`
	Price             float64 `json:"price,omitempty"`
	Total             float64 `json:"total,omitempty"`
	AdditionalCharges float64 `json:"additionalCharges,omitempty"`

	ExpectedDeliveryDate *string `json:"expectedDeliveryDate,omitempty"`
	DeliveredDate        *string `json:"deliveredDate,omitempty"`

	PersonId      *uint `json:"personId,omitempty"`
	MeasurementId *uint `json:"measurementId,omitempty"`
	DressTypeId   *uint `json:"dressTypeId,omitempty"`

	OrderId uint `json:"orderId,omitempty"`

	Files []models.ConfirmFile `json:"files,omitempty"`
}
