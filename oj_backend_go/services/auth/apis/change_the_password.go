package apis

import (
	"github.com/rudrakshsattabhayya/oj_backend_go/config"
	"github.com/rudrakshsattabhayya/oj_backend_go/helpers"
	"github.com/rudrakshsattabhayya/oj_backend_go/models/auth"
	"gorm.io/gorm"
)

type ChangeThePasswordParams struct {
	PerformedByID string `json:"performed_by_id" binding:"required"`
	Password      string `json:"password" binding:"required"`
}

type ChangeThePasswordResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func ChangeThePassword(params ChangeThePasswordParams) (ChangeThePasswordResponse, error) {
	db := config.GetDB()
	tx := db.Begin()
	res, err := ExecuteChangeThePwdTransCode(params, tx)

	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	return res, err
}

func ExecuteChangeThePwdTransCode(params ChangeThePasswordParams, tx *gorm.DB) (ChangeThePasswordResponse, error) {
	// var hashedPassword string

	// if err := tx.Model(&auth.User{}).Select("hashed_password").Where("id = ?", params.PerformedByID).Scan(&hashedPassword).Error; err != nil {
	// 	return res, err
	// }

	// if isCorrect := helpers.VerifyPassword(params.Password, hashedPassword); !isCorrect {
	// 	return res, fmt.Errorf("password is incorrect")
	// }

	if result := tx.Model(&auth.User{}).Where("id = ?", params.PerformedByID).Update("hashed_password", helpers.HashPassword(params.Password)); result.Error != nil {
		return ChangeThePasswordResponse{}, result.Error
	}

	res := GetChangeThePwdResponse()

	return res, nil
}

func GetChangeThePwdResponse() ChangeThePasswordResponse {
	res := ChangeThePasswordResponse{
		Message: "Password is changed!",
		Status:  200,
	}

	return res
}
