package responseModel

type Person struct {
	ID       uint `json:"id,omitempty"`
	IsActive bool `json:"isActive,omitempty"`

	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Gender    string `json:"gender,omitempty"`
	Age       *int   `json:"age,omitempty"`

	CustomerId *uint     `json:"customerId,omitempty"`
	Customer   *Customer `json:"customer,omitempty"`

	Measurements []Measurement `json:"measurements,omitempty"`
}
