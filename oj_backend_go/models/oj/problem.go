package oj

import (
	"time"

	"github.com/google/uuid"
)

type Problem struct {
	ID                  uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Title               string    `gorm:"type:varchar(30);not null" json:"title"`
	ProblemStatement    string    `gorm:"type:varchar(500);not null" json:"problemStatement"`
	VisibleTestCases    string    `gorm:"type:varchar(500);" json:"visibleTestCases,omitempty"`
	VisibleOutputs      string    `gorm:"type:varchar(500);" json:"visibleOutputs,omitempty"`
	HiddenTestCases     string    `gorm:"type:varchar(500);not null" json:"hiddenTestCases"`
	CorrectSolution     string    `gorm:"type:varchar(500);not null" json:"correctSolution"`
	CorrectOutput       string    `gorm:"type:varchar(500);" json:"correctOutput,omitempty"`
	Difficulty          int       `gorm:"default:1" json:"difficulty"`
	AcceptedSubmissions int       `gorm:"default:0" json:"acceptedSubmissions"`
	TotalSubmissions    int       `gorm:"default:0" json:"totalSubmissions"`
	Tags                []Tag     `gorm:"many2many:problem_tags;" json:"tags"`
	CreatedAt           time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
