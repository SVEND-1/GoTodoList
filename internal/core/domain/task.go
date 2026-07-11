package domain

import (
	core_errors "TodoList/internal/core/errors"
	"fmt"
	"time"
)

type Task struct {
	Id          int
	Version     int
	Title       string
	Description *string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time
	UserId      int
}

func NewTask(id int, version int, title string, description *string, completed bool, createAt time.Time, completedAt *time.Time, userId int) Task {
	return Task{
		Id:          id,
		Version:     version,
		Title:       title,
		Description: description,
		Completed:   completed,
		CreatedAt:   createAt,
		CompletedAt: completedAt,
		UserId:      userId,
	}
}

func NewTaskUninitialized(title string, description *string, userId int) Task {
	return Task{
		Id:          UninitializedID,
		Version:     UninitializedVersion,
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
		UserId:      userId,
	}
}

func (t *Task) Validate() error {
	titleLen := len(t.Title)
	if titleLen < 1 || titleLen > 100 {
		return fmt.Errorf("Invalid task title length %d : %w", titleLen, core_errors.ErrInvalidArgument)
	}
	if t.Description != nil {
		descriptionLen := len(*t.Description)
		if descriptionLen < 1 || descriptionLen > 1000 {
			return fmt.Errorf("Invalid task description length %d : %w", descriptionLen, core_errors.ErrInvalidArgument)
		}
	}

	if t.Completed {
		if t.CompletedAt == nil {
			return fmt.Errorf("CompletedAt can`t be `nil` if Completed==true %v: %w", t.CreatedAt, core_errors.ErrInvalidArgument)
		}
		if t.CompletedAt.Before(t.CreatedAt) {
			return fmt.Errorf("completedAt can`t be before createdAt: %w", core_errors.ErrInvalidArgument)
		}
	} else {
		if t.CompletedAt != nil {
			return fmt.Errorf("Completed must be `nil` if Completed==false %v: %w", t.CreatedAt, core_errors.ErrInvalidArgument)
		}
	}

	return nil
}

type TaskPatch struct {
	Title       Nullable[string]
	Description Nullable[string]
	Completed   Nullable[bool]
}

func NewTaskPatch(title Nullable[string], description Nullable[string], completed Nullable[bool]) TaskPatch {
	return TaskPatch{
		Title:       title,
		Description: description,
		Completed:   completed,
	}
}
func (p *TaskPatch) Validate() error {
	if p.Title.Set && p.Title.Value == nil {
		return fmt.Errorf("Title.Value can not be NULL: %w", core_errors.ErrInvalidArgument)
	}
	if p.Completed.Set && p.Completed.Value == nil {
		return fmt.Errorf("Completed.Value can not be NULL: %w", core_errors.ErrInvalidArgument)
	}
	return nil
}

func (t *Task) ApplyPatch(patch TaskPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate task patch %w", err)
	}

	tmp := *t

	if patch.Title.Set {
		tmp.Title = *patch.Title.Value
	}
	if patch.Description.Set {
		tmp.Description = patch.Description.Value
	}
	if patch.Completed.Set {
		tmp.Completed = *patch.Completed.Value

		if tmp.Completed {
			completedAt := time.Now()
			tmp.CompletedAt = &completedAt
		} else {
			tmp.CompletedAt = nil
		}
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate task patch %w", err)
	}

	*t = tmp
	return nil
}
