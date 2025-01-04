package apis

import (
	"github.com/rudrakshsattabhayya/oj_backend_go/config"
	"github.com/rudrakshsattabhayya/oj_backend_go/models/auth"
	"gorm.io/gorm"
)

func InitData() error {
	db := config.GetDB()
	tx := db.Begin()

	err := InitDataTransCode(tx)

	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	return err
}

func InitDataTransCode(tx *gorm.DB) error {
	count, err := CheckUsersCount(tx)
	if count >= 6 || err != nil {
		return err
	}

	if err := CreateFakeUsers(tx); err != nil {
		return err
	}

	if err := CreateFakeProblems(tx); err != nil {
		return err
	}

	return nil
}

func CheckUsersCount(tx *gorm.DB) (int64, error) {
	var count int64
	if err := tx.Model(&auth.User{}).Count(&count).Error; err != nil {
		return -1, err
	}
	return count, nil
}

func CreateFakeUsers(tx *gorm.DB) error {
	users := []auth.User{
		{Email: "dummy1@gmail.com", Username: "Monica", LeaderBoardScore: 20},
		{Email: "dummy2@gmail.com", Username: "Chandler", LeaderBoardScore: 15},
		{Email: "dummy3@gmail.com", Username: "Joey", LeaderBoardScore: 10},
		{Email: "dummy4@gmail.com", Username: "Ross", LeaderBoardScore: 5},
		{Email: "dummy5@gmail.com", Username: "Phoebe", LeaderBoardScore: 5},
		{Email: "dummy6@gmail.com", Username: "Rachel", LeaderBoardScore: 1},
	}

	if err := tx.Create(&users).Error; err != nil {
		return err
	}

	return nil
}

func CreateFakeProblems(tx *gorm.DB) error {
	return nil
}
