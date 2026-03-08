package entities

type Vendor struct {
	*Model `mapstructure:",squash"`

	Name          string `json:"name" gorm:"not null"`
	ContactPerson string `json:"contactPerson"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	Address       string `json:"address" gorm:"type:text"`
	PaymentTerms  string `json:"paymentTerms"`
	Notes         string `json:"notes" gorm:"type:text"`

	// Relations
	Purchases []Purchase `gorm:"foreignKey:VendorId" json:"purchases,omitempty"`
}

func (Vendor) TableNameForQuery() string {
	return "\"stich\".\"Vendors\" E"
}
