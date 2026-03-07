package responseModel

import "time"

type Order struct {
	ID       uint `json:"id,omitempty"`
	IsActive bool `json:"isActive,omitempty"`

	Status string `json:"status,omitempty"`

	Notes string `json:"notes,omitempty"`

	AdditionalCharges float64 `json:"additionalCharges,omitempty"`

	ExpectedDeliveryDate *time.Time `json:"expectedDeliveryDate,omitempty"`
	DeliveredDate        *time.Time `json:"deliveredDate,omitempty"`

	CustomerId   *uint     `json:"customerId,omitempty"`
	Customer     *Customer `json:"customer,omitempty"`
	CustomerName string    `json:"customerName,omitempty"` // first_name + last_name

	OrderTakenById *uint  `json:"orderTakenById,omitempty"`
	OrderTakenBy   string `json:"orderTakenBy,omitempty"` // first_name + last_name

	OrderQuantity int     `json:"orderQuantity,omitempty"` // sum of quantity from order items
	OrderValue    float64 `json:"orderValue,omitempty"`    // sum of total from order items

	AuditFields `json:"auditFields,omitempty"`

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

	ExpectedDeliveryDate *time.Time `json:"expectedDeliveryDate,omitempty"`
	DeliveredDate        *time.Time `json:"deliveredDate,omitempty"`

	PersonId      *uint        `json:"personId,omitempty"`
	Person        *Person      `json:"person,omitempty"`
	MeasurementId *uint        `json:"measurementId,omitempty"`
	Measurement   *Measurement `json:"measurement,omitempty"`

	OrderId uint   `json:"orderId,omitempty"`
	Order   *Order `json:"order,omitempty"`

	Files []EntityDocument `json:"files,omitempty"`

	AuditFields `json:"auditFields,omitempty"`
}
