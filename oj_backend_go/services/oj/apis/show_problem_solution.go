package apis

import (
	"github.com/google/uuid"
	"github.com/rudrakshsattabhayya/oj_backend_go/config"
	"github.com/rudrakshsattabhayya/oj_backend_go/models/oj"
	"gorm.io/gorm"
)

type ShowProblemSolutionParams struct {
	PerformedByID string `json:"performedById" binding:"required"`
	QuestionId    string `form:"questionId" binding:"required"`
}

type ShowProblemSolutionResponse struct {
	Solution string `json:"solution"`
	Status   string `json:"status"`
}

func ShowProblemSolution(params ShowProblemSolutionParams) (ShowProblemSolutionResponse, error) {
	db := config.GetDB()
	tx := db.Begin()

	res, err := ShowProblemSolutionTransCode(params, tx)

	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	return res, err
}

func ShowProblemSolutionTransCode(params ShowProblemSolutionParams, tx *gorm.DB) (ShowProblemSolutionResponse, error) {
	var solution string

	UserID, err := uuid.Parse(params.PerformedByID)
	if err != nil {
		return ShowProblemSolutionResponse{Status: "BAD_REQUEST"}, err
	}

	if err := tx.Model(&oj.Problem{}).Where("id = ?", params.QuestionId).Pluck("correct_solution", &solution).Error; err != nil {
		return ShowProblemSolutionResponse{Status: "BAD_REQUEST"}, err
	}

	res := ShowProblemSolutionResponse{Solution: solution, Status: "success"}

	var count int64
	if err := tx.Model(&oj.ProblemId{}).Where("problem_id = ? AND user_id = ?", params.QuestionId, params.PerformedByID).Count(&count).Error; err != nil {
		return ShowProblemSolutionResponse{Status: "BAD_REQUEST"}, err
	}

	if count == 0 {
		if err := tx.Create(&oj.ProblemId{ProblemID: params.QuestionId, UserID: UserID}).Error; err != nil {
			return ShowProblemSolutionResponse{Status: "BAD_REQUEST"}, err
		}
	}

	return res, nil
}
