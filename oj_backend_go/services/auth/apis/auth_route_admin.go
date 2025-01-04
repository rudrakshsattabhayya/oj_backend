package apis

import (
	"fmt"

	"github.com/rudrakshsattabhayya/oj_backend_go/config"
	"github.com/rudrakshsattabhayya/oj_backend_go/models/auth"
	"gorm.io/gorm"
)

type AuthenticateRouteAdminParams struct {
	PerformedByID string `json:"performed_by_id" binding:"required"`
}

type AuthenticateRouteAdminResponse struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	ProfilePic string `json:"profilePic"`
	Username   string `json:"username"`
	Status     int    `json:"status"`
}

func AuthenticateRouteAdmin(params AuthenticateRouteAdminParams) (AuthenticateRouteAdminResponse, error) {
	db := config.GetDB()

	res, err := ExecuteAuthRouteAdmin(params, db)

	return res, err
}

func ExecuteAuthRouteAdmin(params AuthenticateRouteAdminParams, db *gorm.DB) (AuthenticateRouteAdminResponse, error) {
	var res AuthenticateRouteAdminResponse

	user := auth.User{}

	if err := db.First(&user, "id = ?", params.PerformedByID); err.Error != nil {
		return res, err.Error
	}

	if !user.IsAdmin {
		return res, fmt.Errorf("unauthorized")
	}

	res = GetAuthRouteAdminResp(res, user)

	return res, nil
}

func GetAuthRouteAdminResp(res AuthenticateRouteAdminResponse, user auth.User) AuthenticateRouteAdminResponse {
	res.Name = user.Name
	res.Email = user.Email
	res.ProfilePic = user.ProfilePic
	res.Username = user.Username
	res.Status = 200

	return res
}
