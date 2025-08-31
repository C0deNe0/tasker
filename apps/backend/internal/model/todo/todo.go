package todo

import (
	"time"

	"github.com/google/uuid"
	"github.com/sriniously/go-tasker/internal/model"
	"github.com/sriniously/go-tasker/internal/model/category"
	"github.com/sriniously/go-tasker/internal/model/comment"
)

type Status string

const (
	StatusDraft     Status = "draft"
	StatusActive    Status = "active"
	StatusCompleted Status = "completed"
	StatusArchived  Status = "archived"
)

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

type Todo struct {
	model.Base
	UserID       string      `json:"userId" db:"user_id"`
	Title        string      `json:"title" db:"title"`
	Description  *string     `json:"description" db:"description"`
	Status       Status      `json:"status" db:"status"`
	Priority     Priority    `json:"priority" db:"priority"`
	DueDate      *time.Time  `json:"dueDate" db:"due_date"`
	CompletedAt  *time.Time  `json:"completedAt" db:"completed_at"`
	ParentTodoID *uuid.UUID  `json:"parentTodoId" db:"parent_todo_id"`
	CategoryID   *uuid.UUIDs `json:"categoryId" db:"category_id"`
	MetaData     *MetaData   `json:"metaData" db:"metadata"`
	SortOrder    int         `json:"sortOrder" db:"sort_order"`
}

type MetaData struct {
	Tags       []string `json:"tags"`
	Reminder   *string  `json:"reminder"`
	Color      *string  `json:"color"`
	Difficulty *int     `json:"difficulty"`
}

type PopulatedTodo struct {
	Todo
	Category *category.Category `json:"category" db:"category"`
	Children []Todo             `json:"children" db:"children"`
	Comments []comment.Comment  `json:"comments" db:"comments"`
}
