package models

import (
	"errors"
	"time"
)

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
}

// AllowedStatuses defines the valid status options for tasks
var AllowedStatuses = []string{"pending", "in-progress", "completed"}

func (t *Task) Validate() error {
	if t.ID == "" {
		return errors.New("task ID cannot be empty")
	}
	if t.Title == "" {
		return errors.New("task title cannot be empty")
	}
	if t.DueDate.IsZero() {
		return errors.New("task due date cannot be empty")
	}
	if t.Status == "" {
		return errors.New("task status cannot be empty")
	}
	if !t.isValidStatus() {
		return errors.New("invalid task status. Allowed statuses are: pending, in-progress, completed")
	}
	return nil
}

// isValidStatus checks if the task status is within the allowed statuses
func (t *Task) isValidStatus() bool {
	for _, allowedStatus := range AllowedStatuses {
		if t.Status == allowedStatus {
			return true
		}
	}
	return false
}
