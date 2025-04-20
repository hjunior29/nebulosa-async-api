package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Default struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
