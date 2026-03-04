package responseModel

type Enquiry struct {
	ID         uint      `json:"id,omitempty"`
	IsActive   bool      `json:"isActive,omitempty"`
	Subject    string    `json:"subject,omitempty"`
	Notes      string    `json:"notes,omitempty"`
	Status     string    `json:"status,omitempty"`
	CustomerId *uint     `json:"customerId,omitempty"`
	Customer   *Customer `json:"customer,omitempty"`

	Source              string `json:"source,omitempty"`
	ReferredBy          string `json:"referredBy,omitempty"`
	ReferrerPhoneNumber string `json:"referrerPhoneNumber,omitempty"`

	AuditFields `json:"auditFields,omitempty"`
}
