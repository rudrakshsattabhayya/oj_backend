package apis

import (
	"github.com/rudrakshsattabhayya/oj_backend_go/config"
	"github.com/rudrakshsattabhayya/oj_backend_go/models/auth"
	"github.com/rudrakshsattabhayya/oj_backend_go/models/oj"
	"gorm.io/gorm"
)

type UpdateVerdictParams struct {
	Verdict bool   `form:"verdict" binding:"required"`
	TaskID  string `form:"task_id" binding:"required"`
	Reason  string `form:"reason"`
}

type UpdateVerdictResponse struct {
	Status string `json:"status"`
}

func UpdateVerdict(params UpdateVerdictParams) (UpdateVerdictResponse, error) {
	db := config.GetDB()
	tx := db.Begin()

	res, err := UpdateVerdictTransCode(params, tx)

	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	return res, err
}

func UpdateVerdictTransCode(params UpdateVerdictParams, tx *gorm.DB) (UpdateVerdictResponse, error) {
	submission, err := GetSubmissionObject(params, tx)
	if err != nil {
		return UpdateVerdictResponse{Status: "Error"}, err
	}

	submission.Verdict = params.Verdict
	submission.Reason = params.Reason
	submission.Status = "Processed"

	if err := tx.Save(&submission).Error; err != nil {
		return UpdateVerdictResponse{Status: "Error"}, err
	}

	UpdateSubmissionStats(submission, tx)

	return UpdateVerdictResponse{Status: "Sucess"}, nil
}

func GetSubmissionObject(params UpdateVerdictParams, tx *gorm.DB) (oj.Submission, error) {
	var submission oj.Submission
	if err := tx.Where("request_id = ?", params.TaskID).First(&submission).Error; err != nil {
		return submission, err
	}

	return submission, nil
}

func UpdateSubmissionStats(submission oj.Submission, tx *gorm.DB) {
	var problem oj.Problem
	var user auth.User

	tx.Where("id = ?", submission.ProblemID).First(&problem)
	tx.Where("id = ?", submission.UserID).First(&user)

	user.TotalSubmissions += 1
	problem.TotalSubmissions += 1

	if submission.Verdict{
		var solutionViewed int64
		tx.Model(&oj.ProblemId{}).Where("problem_id = ? AND user_id = ?", problem.ID, user.ID).Count(&solutionViewed)

		var submissionsCount int64
		tx.Model(&oj.Submission{}).Where("user_id = ? AND problem_id = ?", user.ID, problem.ID).Count(&submissionsCount)

		if solutionViewed == 0 && submissionsCount == 1 {
			proofOfSolved := oj.ProblemId{
				ProblemID: problem.ID.String(),
				UserID:    user.ID,
			}
			tx.Create(&proofOfSolved)

			user.LeaderBoardScore += problem.Difficulty
		}

		problem.AcceptedSubmissions += 1
		user.AcceptedSubmissions += 1
	}

	tx.Save(&user)
	tx.Save(&problem)
}
