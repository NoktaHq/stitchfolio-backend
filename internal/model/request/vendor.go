package requestModel

type Vendor struct {
	ID            uint   `json:"id,omitempty"`
	IsActive      bool   `json:"isActive,omitempty"`
	Name          string `json:"name,omitempty"`
	ContactPerson string `json:"contactPerson,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Email         string `json:"email,omitempty"`
	Address       string `json:"address,omitempty"`
	PaymentTerms  string `json:"paymentTerms,omitempty"`
	Notes         string `json:"notes,omitempty"`
}
