package requestModel

type Task struct {
	ID           uint    `json:"id,omitempty"`
	IsActive     *bool   `json:"isActive,omitempty"`
	Title        string  `json:"title"`
	Description  *string `json:"description,omitempty"`
	Status       string  `json:"status,omitempty"` // PENDING, IN_PROGRESS, COMPLETED, CANCELLED
	Priority     *int    `json:"priority,omitempty"`
	DueDate      *string `json:"dueDate,omitempty"`
	ReminderDate *string `json:"reminderDate,omitempty"`
	CompletedAt  *string `json:"completedAt,omitempty"`
	AssignedToId *uint   `json:"assignedToId,omitempty"`
}
