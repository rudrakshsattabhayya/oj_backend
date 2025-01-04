package apis

import (
	"github.com/rudrakshsattabhayya/oj_backend_go/config"
	"github.com/rudrakshsattabhayya/oj_backend_go/models/auth"
	"gorm.io/gorm"
)

type AuthenticateRouteParams struct {
	PerformedByID string `json:"performed_by_id" binding:"required"`
}

type AuthenticateRouteResponse struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	ProfilePic string `json:"profilePic"`
	Username   string `json:"username"`
	Status     int    `json:"status"`
}

func AuthenticateRoute(params AuthenticateRouteParams) (AuthenticateRouteResponse, error) {
	db := config.GetDB()

	res, err := ExecuteAuthRoute(params, db)

	return res, err
}

func ExecuteAuthRoute(params AuthenticateRouteParams, db *gorm.DB) (AuthenticateRouteResponse, error) {
	var res AuthenticateRouteResponse

	user := auth.User{}

	if err := db.First(&user, "id = ?", params.PerformedByID); err.Error != nil {
		return res, err.Error
	}

	res = GetAuthRouteResp(res, user)

	return res, nil
}

func GetAuthRouteResp(res AuthenticateRouteResponse, user auth.User) AuthenticateRouteResponse {
	res.Name = user.Name
	res.Email = user.Email
	res.ProfilePic = user.ProfilePic
	res.Username = user.Username
	res.Status = 200

	return res
}
