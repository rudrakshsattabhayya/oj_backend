package apis

import (
	"fmt"
	"time"

	"github.com/rudrakshsattabhayya/oj_backend_go/config"
	"gorm.io/gorm"
)

type ListSubmissionsParams struct {
	PerformedByID string `json:"performedById" binding:"required"`
}
type ListSubmissionsResponse struct {
	Submissions []struct {
		ID      string    `json:"id"`
		User    string `json:"user"`
		Problem struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		} `json:"problem"`
		Code      string    `json:"code"`
		Time      time.Time `json:"time"`
		Verdict   bool      `json:"verdict"`
		Reason    string    `json:"reason"`
		RequestID string    `json:"request_id"`
		Status    string    `json:"status"`
	} `json:"submissions"`
	Status int `json:"status"`
}

func ListSubmissions(params ListSubmissionsParams) (ListSubmissionsResponse, error) {
	db := config.GetDB()

	res, err := GetListSubmissionsResponse(params, db)
	if err != nil {
		return ListSubmissionsResponse{}, err
	}

	return res, nil
}

func GetListSubmissionsResponse(params ListSubmissionsParams, db *gorm.DB) (ListSubmissionsResponse, error) {
	res := ListSubmissionsResponse{
		Status: 200,
	}
	fmt.Println("params.PerformedByID ", params.PerformedByID)

	if err := db.Table("submissions s").
		Select("s.id, u.name AS user, jsonb_build_object('id', p.id, 'title', p.title) AS problem, s.code, s.time, s.verdict, s.reason, s.request_id, s.status").
		Joins("INNER JOIN problems p ON s.problem_id = p.id").
		Joins("INNER JOIN users u ON s.user_id = u.id").
		Where("s.user_id = ?", params.PerformedByID).
		Scan(&res.Submissions).Error; err != nil {
		return ListSubmissionsResponse{Status: 400}, err
	}

	return res, nil
}
