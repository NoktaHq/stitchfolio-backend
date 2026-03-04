package entities

type Category struct {
	*Model `mapstructure:",squash"`

	Name string `json:"name" gorm:"not null"`

	// Relations
	Products []Product `gorm:"foreignKey:CategoryId;constraint:OnDelete:SET NULL" json:"products,omitempty"`
}

func (Category) TableNameForQuery() string {
	return "\"stich\".\"Categories\" E"
}
