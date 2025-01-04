package helpers

import (
	"fmt"

	"github.com/rudrakshsattabhayya/oj_backend_go/config"
	"github.com/rudrakshsattabhayya/oj_backend_go/models/auth"
)

func AuthenticateWithToken(performed_by_id string, token string) error {
	db := config.GetDB()

	var count int64

	if err := db.Model(&auth.User{}).Where("id = ? AND token = ?", performed_by_id, token).Count(&count).Error; err != nil{
		return fmt.Errorf("token is invalid")
	}

	return nil
}

func AuthenticateWithPassword(performed_by_id string, token string) error {
	return nil
}
