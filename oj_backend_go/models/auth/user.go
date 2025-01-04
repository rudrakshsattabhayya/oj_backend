package auth

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                  uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Email               string    `gorm:"type:varchar(100);not null;unique" json:"email"`
	HashedPassword      string    `gorm:"type:varchar(1000);" json:"hashedPassword,omitempty"`
	Name                string    `gorm:"type:varchar(50);not null" json:"name"`
	Username            string    `gorm:"type:varchar(20);not null;unique" json:"username"`
	ProfilePic          string    `gorm:"type:varchar(500);not null" json:"profilePic"`
	IsAdmin             bool      `gorm:"default:false" json:"isAdmin"`
	Token               string    `gorm:"type:varchar(1000);default:''" json:"token"`
	LeaderBoardScore    int       `gorm:"default:0" json:"leaderBoardScore"`
	TotalSubmissions    int       `gorm:"default:0" json:"totalSubmissions"`
	AcceptedSubmissions int       `gorm:"default:0" json:"acceptedSubmissions"`
	CreatedAt           time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
