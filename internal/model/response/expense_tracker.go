package responseModel

import "time"

type ExpenseTracker struct {
	ID       uint `json:"id,omitempty"`
	IsActive bool `json:"isActive,omitempty"`

	PurchaseDate *time.Time `json:"purchaseDate,omitempty"`
	BillNumber   string     `json:"billNumber,omitempty"`
	CompanyName  string     `json:"companyName,omitempty"`
	Material     string     `json:"material,omitempty"`
	Price        float64    `json:"price,omitempty"`
	Balance      float64    `json:"balance,omitempty"`
	Location     *string    `json:"location,omitempty"`
	Notes        *string    `json:"notes,omitempty"`

	ExpenseDetails []ExpenseDetail `json:"expenseDetails,omitempty"`
	Files          []EntityDocument `json:"files,omitempty"`

	AuditFields `json:"auditFields,omitempty"`
}

type ExpenseDetail struct {
	ID        uint    `json:"id,omitempty"`
	IsActive  bool    `json:"isActive,omitempty"`
	Source    string  `json:"source,omitempty"`
	Price     float64 `json:"price,omitempty"`
	ExpenseId uint    `json:"expenseId,omitempty"`

	AuditFields `json:"auditFields,omitempty"`
}
