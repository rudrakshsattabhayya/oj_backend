package oj

import (
	"github.com/google/uuid"
)

type ProblemTag struct {
	ProblemID uuid.UUID `gorm:"type:uuid;not null;primaryKey"`
	TagID     uuid.UUID `gorm:"type:uuid;not null;primaryKey"`
}

func (ProblemTag) TableName() string {
	return "problem_tags"
}
