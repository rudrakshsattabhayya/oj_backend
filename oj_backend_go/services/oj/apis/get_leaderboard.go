package apis

import (
	"github.com/rudrakshsattabhayya/oj_backend_go/config"
	"github.com/rudrakshsattabhayya/oj_backend_go/models/auth"
	"gorm.io/gorm"
)

type GetLeaderboardParams struct {
}

type GetLeaderboardResponse struct {
	Response []struct {
		ID               string `json:"id"`
		Username         string `json:"username"`
		LeaderBoardScore int `json:"leaderBoardScore"`
	} `json:"response"`
	Status int `json:"status"`
}

func GetLeaderboard(params GetLeaderboardParams) (GetLeaderboardResponse, error) {
	db := config.GetDB()

	res, err := GetLeaderboardData(db)
	if err != nil {
		return GetLeaderboardResponse{Status: 400}, err
	}

	return res, nil
}

func GetLeaderboardData(db *gorm.DB) (GetLeaderboardResponse, error) {
	res := GetLeaderboardResponse{Status: 200}

	if err := db.Model(&auth.User{}).Select("id, username, leader_board_score").Scan(&res.Response).Error; err != nil {
		return GetLeaderboardResponse{Status: 400}, err
	}

	return res, nil
}
