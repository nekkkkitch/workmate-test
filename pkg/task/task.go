package task

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Status string

var (
	Created    Status = "CREATED"
	InProgress Status = "INPROGRESS"
	Finished   Status = "FINISHED"
	Cancelled  Status = "Cancelled"
)

type Task struct {
	ID             uuid.UUID          `json:"id"`
	Status         Status             `json:"status"`
	CreationDate   time.Time          `json:"creation_date"`
	TimeSpent      time.Duration      `json:"time_spent"`
	Context        context.Context    `json:"-"`
	CancelFunction context.CancelFunc `json:"-"`
}
