package apis

import (
	"github.com/rudrakshsattabhayya/oj_backend_go/config"
	"github.com/rudrakshsattabhayya/oj_backend_go/models/auth"
	"gorm.io/gorm"
)

type ChangeUsernameParams struct {
	PerformedByID string `json:"performed_by_id" binding:"required"`
	Username      string `json:"newUserName" binding:"required"`
}

type ChangeUsernameResponse struct {
	Username string `json:"username"`
	Message  string `json:"message"`
	Status   int    `json:"status"`
}

func ChangeUsername(params ChangeUsernameParams) (ChangeUsernameResponse, error) {
	db := config.GetDB()
	tx := db.Begin()
	res, err := ExecuteChangeUsernameTransCode(params, tx)

	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	return res, err
}

func ExecuteChangeUsernameTransCode(params ChangeUsernameParams, tx *gorm.DB) (ChangeUsernameResponse, error) {
	var res ChangeUsernameResponse

	if result := tx.Model(&auth.User{}).Where("id = ?", params.PerformedByID).Update("username", params.Username); result.Error != nil {
		return res, result.Error
	}

	res = GetChangeUsernameResponse(res, params.Username)

	return res, nil
}

func GetChangeUsernameResponse(res ChangeUsernameResponse, username string) ChangeUsernameResponse {
	res.Message = "Username is changed!"
	res.Status = 200
	res.Username = username

	return res
}
