package oj

import (
	"time"

	"github.com/google/uuid"
	"github.com/rudrakshsattabhayya/oj_backend_go/models/auth"
)

type ProblemId struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ProblemID string     `gorm:"type:varchar(500);not null" json:"problemId"`
	UserID uuid.UUID `gorm:"type:uuid" json:"userId"`
	User      *auth.User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"createdAt"`
}
