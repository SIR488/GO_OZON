package domain

import (
	"errors"
	"time"
)

type Task struct {
	ID       ID
	Title    string
	Status   Status
	Deadline time.Time
}
type ID uint

type Status string

const (
	StatusPlanned    Status = "planned"
	StatusInProgress Status = "in_progress"
	StatusCanceled   Status = "canceled"
	StatusDone       Status = "done"
)

var ErrNotFound error = errors.New("not found")
