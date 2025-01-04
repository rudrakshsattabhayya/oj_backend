package apis

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"

	"net/http"
	"net/url"

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

	submission, err := CreateSubmission(params, problem, tx)
	if err != nil {
		return SubmitProblemResponse{Status: 400, Message: "Error: "}, err
	}

	task_id, err := EvaluateSubmission(submission, problem, tx)
	if err != nil {
		return SubmitProblemResponse{}, err
	}

	if err := UpdateTaskID(submission, task_id, tx); err != nil {
		return SubmitProblemResponse{}, err
	}

	return SubmitProblemResponse{Verdict: "Queued", Status: 200, Message: "Successful submission!", TaskID: task_id}, nil
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
		Code:      codeurl,
		UserID:    UserID,
		ProblemID: problem.ID,
		Verdict:   false,
	}

	if err := tx.Create(&submission).Error; err != nil {
		return oj.Submission{}, err
	}

	return submission, nil
}

func EvaluateSubmission(submission oj.Submission, problem oj.Problem, tx *gorm.DB) (string, error) {
	baseURL := os.Getenv("BACKEND_DJANGO_EVALUATION_SERVER_URL") + "/get-verdict"
	params := url.Values{}
	params.Add("code", submission.Code)
	params.Add("inputs", problem.HiddenTestCases)
	params.Add("correctOutputs", problem.CorrectOutput)
	params.Add("password", os.Getenv("BACKEND_DJANGO_PASSWORD"))

	resp, err := http.PostForm(baseURL, params)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatal(err)
	}

	task_id := result["task_id"].(string)

	return task_id, nil
}

func UpdateTaskID(submission oj.Submission, task_id string, tx *gorm.DB) error {
	TaskID, err := uuid.Parse(task_id)
	if err != nil {
		return err
	}

	submission.RequestID = TaskID
	if err := tx.Save(&submission).Error; err != nil {
		return err
	}

	return nil
}
