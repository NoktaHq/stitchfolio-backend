package responseModel

import "time"

type EnquiryHistory struct {
	ID              uint      `json:"id,omitempty"`
	IsActive        bool      `json:"isActive,omitempty"`
	Status          *string   `json:"status,omitempty"`
	EmployeeComment string    `json:"employeeComment,omitempty"`
	CustomerComment string    `json:"customerComment,omitempty"`
	VisitingDate    *string   `json:"visitingDate,omitempty"`
	CallBackDate    *string   `json:"callBackDate,omitempty"`
	EnquiryDate     string    `json:"enquiryDate,omitempty"`
	ResponseStatus  string    `json:"responseStatus,omitempty"`
	EnquiryId       uint      `json:"enquiryId,omitempty"`
	EmployeeId      uint      `json:"employeeId,omitempty"`
	Employee        *User     `json:"employee,omitempty"`
	PerformedAt     time.Time `json:"performedAt,omitempty"`
	PerformedById   uint      `json:"performedById,omitempty"`
	PerformedBy     *User     `json:"performedBy,omitempty"`

	AuditFields `json:"auditFields"`
}
