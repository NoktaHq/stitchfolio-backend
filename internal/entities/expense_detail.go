package entities

type ExpenseDetail struct {
	*Model `mapstructure:",squash"`

	Source    string  `json:"source"`
	Price     float64 `json:"price"`
	ExpenseId uint    `json:"expenseId"`

	Expense *Expense `gorm:"foreignKey:ExpenseId" json:"expense,omitempty"`
}

func (ExpenseDetail) TableName() string {
	return TableNameWithSchema("ExpenseDetails")
}

func (ExpenseDetail) TableNameForQuery() string {
	return "\"stich\".\"ExpenseDetails\" E"
}
