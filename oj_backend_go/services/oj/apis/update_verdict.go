package apis

import (
	"github.com/rudrakshsattabhayya/oj_backend_go/config"
	"github.com/rudrakshsattabhayya/oj_backend_go/models/oj"
	"gorm.io/gorm"
)

type UpdateVerdictParams struct {
	Verdict bool `form:"verdict" binding:"required"`
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

	return UpdateVerdictResponse{Status: "Sucess"}, nil
}

func GetSubmissionObject(params UpdateVerdictParams, tx *gorm.DB) (oj.Submission, error) {
	var submission oj.Submission
	if err := tx.Where("request_id = ?", params.TaskID).First(&submission).Error; err != nil {
		return submission, err
	}

	return submission, nil
}
