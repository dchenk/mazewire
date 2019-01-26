package tasks

import (
	"context"
	"net/http"
)

// A TaskName specifies a kind of task.
type TaskName string

// New schedules a task to be executed asynchronously and retried as needed.
func New(ctx context.Context, taskName TaskName, data []byte) error {
	return nil
}

// A TaskHandler is able to handle a particular named task and return an HTTP status code.
// Success is indicated by a 200 status code, and any other value indicates that the task must
// be retried at a later time.
type TaskHandler func(ctx context.Context, r *http.Request) int

const (
	// SendEmail is a task to send an email.
	SendEmail TaskName = "send_email"
)
