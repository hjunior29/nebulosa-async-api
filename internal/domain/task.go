package domain

import (
	"time"

	"gorm.io/datatypes"
)

type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusProcessing TaskStatus = "processing"
	StatusSuccess    TaskStatus = "success"
	StatusFailed     TaskStatus = "failed"
)

type Task struct {
	Default
	Endpoint        string         `json:"endpoint"`
	Headers         datatypes.JSON `json:"headers" gorm:"type:jsonb"`
	Method          string         `json:"method"`
	Payload         datatypes.JSON `json:"payload" gorm:"type:jsonb"`
	Type            string         `json:"type"`
	Status          TaskStatus     `json:"status"`
	MaxRetries      int            `json:"maxRetries"`
	Attempts        int            `json:"attempts"`
	ScheduledAt     string         `json:"scheduledAt"`
	ScheduledAtTime time.Time      `json:"scheduledAtTime"`
	LastError       string         `json:"lastError"`
	StatusCode      int            `json:"statusCode"`
}
