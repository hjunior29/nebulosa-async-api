package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusProcessing TaskStatus = "processing"
	StatusSuccess    TaskStatus = "success"
	StatusFailed     TaskStatus = "failed"
)

type Task struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Endpoint    string         `json:"endpoint"`
	Headers     datatypes.JSON `json:"headers" gorm:"type:jsonb"`
	Method      string         `json:"method"`
	Payload     datatypes.JSON `json:"payload" gorm:"type:jsonb"`
	Type        string         `json:"type"`
	Status      TaskStatus     `json:"status"`
	MaxRetries  int            `json:"maxRetries"`
	Attempts    int            `json:"attempts"`
	ScheduledAt time.Time      `json:"scheduledAt"`
	LastError   string         `json:"lastError"`
	StatusCode  int            `json:"statusCode"`
}
