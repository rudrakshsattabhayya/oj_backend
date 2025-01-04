package apis

import (
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/rudrakshsattabhayya/oj_backend_go/config"
	"github.com/rudrakshsattabhayya/oj_backend_go/helpers"
	"github.com/rudrakshsattabhayya/oj_backend_go/models/oj"
	"gorm.io/gorm"
)

type SubmitProblemParams struct {
	PerformedByID string                `form:"performed_by_id" binding:"required"`
	QuestionId    string                `form:"questionId" binding:"required"`
	Code          *multipart.FileHeader `form:"code" binding:"required"`
}

type SubmitProblemResponse struct {
	Verdict string `json:"verdict"`
	TaskID  string `json:"task_id"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func SubmitProblem(params SubmitProblemParams) (SubmitProblemResponse, error) {
	db := config.GetDB()
	tx := db.Begin()

	res, err := SubmitProblemTranCode(params, tx)

	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	return res, err
}

func SubmitProblemTranCode(params SubmitProblemParams, tx *gorm.DB) (SubmitProblemResponse, error) {
	problem, err := ValidateProblem(params.QuestionId, tx)
	if err != nil {
		return SubmitProblemResponse{Status: 400, Message: "Error: "}, err
	}

	_, err = CreateSubmission(params, problem, tx)
	if err != nil {
		return SubmitProblemResponse{Status: 400, Message: "Error: "}, err
	}

	//Evaluate the submission
	//Update submission with task_id

	return SubmitProblemResponse{Verdict: "Queued", Status: 200, Message: "Successful submission!"}, nil
}

func ValidateProblem(questionId string, tx *gorm.DB) (oj.Problem, error) {
	problem := oj.Problem{}
	result := tx.Where("id = ?", questionId).First(&problem)

	if result.Error != nil {
		return oj.Problem{}, result.Error
	}

	return problem, nil
}

func CreateSubmission(params SubmitProblemParams, problem oj.Problem, tx *gorm.DB) (oj.Submission, error) {
	codeurl, err := helpers.UploadFile(params.Code, "submissions")
	if err != nil {
		return oj.Submission{}, err
	}

	UserID, err := uuid.Parse(params.PerformedByID)
	if err != nil {
		return oj.Submission{}, err
	}

	submission := oj.Submission{
		Code:    codeurl,
		UserID:  UserID,
		ProblemID: problem.ID,
		Verdict: false,
	}

	if err := tx.Create(&submission).Error; err != nil {
		return oj.Submission{}, err
	}

	return submission, nil
}
