package oj

import (
	"time"

	"github.com/google/uuid"
	"github.com/rudrakshsattabhayya/oj_backend_go/models/auth"
)

type Submission struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Code      string     `gorm:"type:varchar(500);not null" json:"code"`
	Time      time.Time  `gorm:"autoCreateTime" json:"time"`
	Verdict   bool       `gorm:"default:false" json:"verdict"`
	Reason    string     `gorm:"type:text;" json:"reason,omitempty"`
	UserID    uuid.UUID  `gorm:"type:uuid" json:"userId"`
	User      *auth.User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user"`
	ProblemID uuid.UUID  `gorm:"type:uuid" json:"problemId"`
	Problem   *Problem   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"problem"`
	RequestID uuid.UUID `gorm:"type:uuid;unique;" json:"requestId,omitempty"`
	Status    string     `gorm:"type:varchar(50);default:'Queued';not null" json:"status"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"createdAt"`
}
