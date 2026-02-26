package responseModel

import "time"

type Task struct {
	ID           uint       `json:"id,omitempty"`
	IsActive     bool       `json:"isActive,omitempty"`
	Title        string     `json:"title,omitempty"`
	Description  *string    `json:"description,omitempty"`
	IsCompleted  bool       `json:"isCompleted"`
	Priority     *int       `json:"priority,omitempty"`
	DueDate      *time.Time `json:"dueDate,omitempty"`
	ReminderDate *time.Time `json:"reminderDate,omitempty"`
	CompletedAt  *time.Time `json:"completedAt,omitempty"`
	AssignedToId *uint      `json:"assignedToId,omitempty"`

	AuditFields `json:"auditFields,omitempty"`
}
