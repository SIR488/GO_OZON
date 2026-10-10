package service

import (
	"MyMod/domain"
	"errors"
	"fmt"
	"sort"
	"time"
)

type Repository interface {
	Save(domain.Task) error
	ReadAll() ([]domain.Task, error)
	ReadByID(domain.ID) (domain.Task, error)
	Delete(domain.ID) error
	ChangeTask(domain.ID, *domain.Status, *string, *time.Time) error
}

type MainService struct {
	Repository Repository
}

func (s *MainService) Save(title string, deadline time.Time) error {
	if title == "" {
		return fmt.Errorf("empty title")
	}
	if deadline.Before(time.Now()) {
		return fmt.Errorf("deadline has passed")
	}
	task := domain.Task{
		Title:    title,
		Status:   domain.StatusPlanned,
		Deadline: deadline,
	}
	return s.Repository.Save(task)
}

func (s *MainService) ReadAll() ([]domain.Task, error) {
	tasks, err := s.Repository.ReadAll()
	if err != nil {
		return []domain.Task{}, err
	}
	sortFunc := func(i int, j int) bool {
		if tasks[i].Deadline.Equal(tasks[j].Deadline) {
			return tasks[i].ID < tasks[j].ID
		}
		return tasks[i].Deadline.Before(tasks[j].Deadline)
	}
	sort.Slice(tasks, sortFunc)
	return tasks, nil
}
func (s *MainService) ReadByID(id domain.ID) (domain.Task, error) {
	return s.Repository.ReadByID(id)
}
func (s *MainService) Delete(id domain.ID) error {
	return s.Repository.Delete(id)
}

func (s *MainService) ChangeStatus(id domain.ID, newStatus domain.Status) error {
	task, err := s.Repository.ReadByID(id)
	if err != nil {
		return fmt.Errorf("search task: %w", err)
	}
	switch task.Status {
	case domain.StatusCanceled:
		return errors.New("task already canceled")
	case domain.StatusDone:
		return errors.New("task already done")
	case newStatus:
		return errors.New("task have same status")
	}
	err = s.Repository.ChangeTask(id, &newStatus, nil, nil)
	if err != nil {
		return fmt.Errorf("valid data. %w", err)
	}
	return nil
}
func (s *MainService) ChangeTask(id domain.ID, newTitle *string, newDeadline *time.Time) error {
	if newTitle == nil && newDeadline == nil {
		return errors.New("one parametrs must be not empty")
	}
	task, err := s.Repository.ReadByID(id)

	switch task.Status {
	case domain.StatusCanceled:
		return errors.New("task already canceled")
	case domain.StatusDone:
		return errors.New("task already done")
	}

	if err != nil {
		return fmt.Errorf("search task: %w", err)
	}

	if newDeadline.Before(time.Now()) {
		return fmt.Errorf("deadline has passed")
	}

	err = s.Repository.ChangeTask(id, nil, newTitle, newDeadline)
	if err != nil {
		return fmt.Errorf("valid data. %w", err)
	}
	return nil
}
