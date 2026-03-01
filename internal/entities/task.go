package entities

import "time"

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "PENDING"
	TaskStatusInProgress TaskStatus = "IN_PROGRESS"
	TaskStatusCompleted  TaskStatus = "COMPLETED"
	TaskStatusCancelled  TaskStatus = "CANCELLED"
)

type Task struct {
	*Model `mapstructure:",squash"`

	Title        string     `gorm:"size:255;not null" json:"title"`
	Description  *string    `gorm:"type:text" json:"description,omitempty"`
	IsCompleted  bool       `gorm:"default:false" json:"isCompleted"`
	Status       TaskStatus `gorm:"default:'PENDING';type:text" json:"status"`
	Priority     *int       `json:"priority,omitempty"`
	DueDate      *time.Time `json:"dueDate,omitempty"`
	ReminderDate *time.Time `json:"reminderDate,omitempty"`
	CompletedAt  *time.Time `json:"completedAt,omitempty"`

	AssignedToId *uint `json:"assignedToId,omitempty"`
	AssignedTo   *User `gorm:"foreignKey:AssignedToId" json:"assignedTo,omitempty"`
}

func (Task) TableNameForQuery() string {
	return "\"stich\".\"Tasks\" E"
}
