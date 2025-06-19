package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Task represents a task in the system
type Task struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Title       string    `json:"title" gorm:"not null;size:255"`
	Description string    `json:"description" gorm:"type:text"`
	Status      string    `json:"status" gorm:"not null;default:'pending';size:50"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"not null;default:now()"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (task *Task) BeforeCreate(tx *gorm.DB) error {
	if task.ID == uuid.Nil {
		task.ID = uuid.New()
	}
	return nil
}

// CreateTaskRequest represents the request body for creating a task
type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required" example:"Complete project documentation"`
	Description string `json:"description" example:"Write comprehensive documentation for the API"`
	Status      string `json:"status" example:"pending"`
}

// UpdateTaskRequest represents the request body for updating a task
type UpdateTaskRequest struct {
	Title       string `json:"title" example:"Complete project documentation"`
	Description string `json:"description" example:"Write comprehensive documentation for the API"`
	Status      string `json:"status" example:"completed"`
}

// TaskResponse represents the response body for task operations
type TaskResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TasksResponse represents the response body for listing tasks
type TasksResponse struct {
	Tasks []TaskResponse `json:"tasks"`
	Total int64          `json:"total"`
} 