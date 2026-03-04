package entities

import "time"

type Expense struct {
	*Model `mapstructure:",squash"`

	PurchaseDate *time.Time `json:"purchaseDate,omitempty"`
	BillNumber   string     `json:"billNumber"`
	CompanyName  string     `json:"companyName"`
	Material     string     `json:"material"`
	Price        float64    `json:"price,omitempty"`
	Balance      float64    `json:"balance,omitempty"` // Price - Sum(ExpenseDetail.Price)
	Location     *string    `json:"location,omitempty"`
	Notes        *string    `json:"notes,omitempty"`

	ExpenseDetails []ExpenseDetail `gorm:"foreignKey:ExpenseId" json:"expenseDetails,omitempty"`
}

func (Expense) TableNameForQuery() string {
	return "\"stich\".\"Expenses\" E"
}
