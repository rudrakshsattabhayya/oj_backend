package apis

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/rudrakshsattabhayya/oj_backend_go/config"
	"github.com/rudrakshsattabhayya/oj_backend_go/models/auth"
	"github.com/rudrakshsattabhayya/oj_backend_go/helpers"
	"gorm.io/gorm"
)

type LoginParams struct {
	JwtToken string `json:"jwtToken" binding:"required"`
}

type LoginResponse struct {
	JwtToken   string `json:"jwtToken"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	ProfilePic string `json:"profilePic"`
	Status     int    `json:"status"`
}

func Login(params LoginParams) (LoginResponse, error) {
	db := config.GetDB()
	tx := db.Begin()
	res, err := ExecuteLoginTransactionCode(params, tx)

	if err != nil {
		tx.Rollback()
	}else{
		tx.Commit()
	}

	return res, err
}

func ExecuteLoginTransactionCode(params LoginParams, tx *gorm.DB) (LoginResponse, error) {
	res := LoginResponse{}

	userDetails, err := DecodeJWTWithoutVerification(params.JwtToken)
	if err != nil {
		return res, err
	}

	user, err := FindOrInitializeUser(userDetails, tx)
	if err != nil {
		return res, err
	}

	res = GetLoginResponse(user, res)

	return res, nil
}

func FindOrInitializeUser(userDetails jwt.MapClaims, tx *gorm.DB) (auth.User, error) {
	email := userDetails["email"].(string)

	user := auth.User{}

	if result := tx.First(&user, "email = ?", email); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {

			user = auth.User{
				Email:      email,
				Name:       userDetails["name"].(string),
				ProfilePic: userDetails["picture"].(string),
				Token:      uuid.New().String(),
			}

			if err := tx.Create(&user).Error; err != nil {
				return user, errors.New(err.Error())
			}
		} else {
			return user, result.Error
		}
	}

	return user, nil
}

func DecodeJWTWithoutVerification(tokenString string) (jwt.MapClaims, error) {
	parser := &jwt.Parser{SkipClaimsValidation: true}
	token, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("invalid JWT: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, fmt.Errorf("unable to extract claims from JWT")
}

func GetLoginResponse(user auth.User, res LoginResponse) LoginResponse {
	res.JwtToken = helpers.CreateSessionToken(user)
	res.Name = user.Name
	res.Email = user.Email
	res.ProfilePic = user.ProfilePic
	res.Status = 200

	return res
}
