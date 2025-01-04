package apis

import (
	"fmt"

	"github.com/rudrakshsattabhayya/oj_backend_go/config"
	"github.com/rudrakshsattabhayya/oj_backend_go/helpers"
	"github.com/rudrakshsattabhayya/oj_backend_go/models/auth"
	"gorm.io/gorm"
)

type LoginWithPasswordParams struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginWithPasswordResponse struct {
	JwtToken   string `json:"jwtToken"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	ProfilePic string `json:"profilePic"`
	Status     int    `json:"status"`
}

func LoginWithPassword(params LoginWithPasswordParams) (LoginWithPasswordResponse, error) {
	db := config.GetDB()
	tx := db.Begin()
	res, err := ExecuteLoginWithPassTransactionCode(params, tx)

	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	return res, err
}

func ExecuteLoginWithPassTransactionCode(params LoginWithPasswordParams, tx *gorm.DB) (LoginWithPasswordResponse, error) {
	res := LoginWithPasswordResponse{}

	user, err := GetUserWithEmail(params.Email, tx)
	if err != nil {
		return res, err
	}

	if isCorrect := helpers.VerifyPassword(params.Password, user.HashedPassword); !isCorrect {
		return res, fmt.Errorf("password is incorrect")
	}

	res = GetLoginWithPasswordResponse(user, res)

	return res, nil
}

func GetUserWithEmail(email string, tx *gorm.DB) (auth.User, error) {
	user := auth.User{}

	if result := tx.First(&user, "email = ?", email); result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

func GetLoginWithPasswordResponse(user auth.User, res LoginWithPasswordResponse) LoginWithPasswordResponse {
	res.JwtToken = helpers.CreateSessionToken(user)
	res.Name = user.Name
	res.Email = user.Email
	res.ProfilePic = user.ProfilePic
	res.Status = 200

	return res
}
