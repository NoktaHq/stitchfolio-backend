package requestModel

type ExpenseTracker struct {
	ID       uint  `json:"id,omitempty"`
	IsActive *bool `json:"isActive,omitempty"`

	PurchaseDate *string `json:"purchaseDate,omitempty"`
	BillNumber   string  `json:"billNumber"`
	CompanyName  string  `json:"companyName"`
	Material     string  `json:"material"`
	Price        float64 `json:"price,omitempty"`
	Balance      float64 `json:"balance,omitempty"`
	Location     *string `json:"location,omitempty"`
	Notes        *string `json:"notes,omitempty"`

	ExpenseDetails []ExpenseDetail `json:"expenseDetails,omitempty"`
}

type ExpenseDetail struct {
	ID         uint    `json:"id,omitempty"`
	IsActive   *bool   `json:"isActive,omitempty"`
	Source     string  `json:"source,omitempty"`
	Price      float64 `json:"price"`
	ExpenseId  uint    `json:"expenseId,omitempty"`
}
