package apis

import (
	"github.com/lib/pq"
	"github.com/rudrakshsattabhayya/oj_backend_go/config"
	"github.com/rudrakshsattabhayya/oj_backend_go/models/oj"
	"gorm.io/gorm"
)

type ShowProblemParams struct {
	PerformedByID string `json:"performedById" binding:"required"`
	QuestionId    string `form:"questionId" binding:"required"`
}

type ShowProblemResponse struct {
	Response struct {
		ProblemsData    []ShowProblemProblemsData `json:"problemsData"`
		UserSubmissions []ShowProblemUserData     `json:"userSubmissions"`
	} `json:"response"`
	Status int `json:"status"`
}

type ShowProblemProblemsData struct {
	ID                  string         `json:"id"`
	Tags                pq.StringArray `gorm:"type:text[]" json:"tags"`
	Title               string         `json:"title"`
	ProblemStatement    string         `json:"problemStatement"`
	VisibleTestCases    string         `json:"visibleTestCases"`
	VisibleOutputs      string         `json:"visibleOutputs"`
	Difficulty          int            `json:"difficulty"`
	AcceptedSubmissions int            `json:"acceptedSubmissions"`
	TotalSubmissions    int            `json:"totalSubmissions"`
}

type ShowProblemUserData struct {
	Time    string `json:"time"`
	Verdict bool   `json:"verdict"`
	Code    string `json:"code"`
	Status  string `json:"status"`
}

func ShowProblem(params ShowProblemParams) (ShowProblemResponse, error) {
	db := config.GetDB()

	res, err := GetShowProblem(params, db)
	if err != nil {
		return ShowProblemResponse{Status: 400}, err
	}

	return res, nil
}

func GetShowProblem(params ShowProblemParams, db *gorm.DB) (ShowProblemResponse, error) {
	res := ShowProblemResponse{Status: 200}

	if err := db.Model(&oj.Problem{}).
		Select(`
			problems.id, 
			problems.title, 
			problems.problem_statement, 
			problems.visible_test_cases, 
			problems.visible_outputs, 
			problems.difficulty, 
			problems.accepted_submissions, 
			problems.total_submissions,
			ARRAY_AGG(COALESCE(t.name, '')) AS tags`).
		Joins("LEFT JOIN problem_tags pt ON pt.problem_id = problems.id").
		Joins("LEFT JOIN tags t ON t.id = pt.tag_id").
		Where("problems.id = ?", params.QuestionId).
		Group("problems.id").
		Find(&res.Response.ProblemsData).Error; err != nil {
		return ShowProblemResponse{Status: 400}, err
	}

	if err := db.Table("submissions s").
		Select("s.id, s.code, s.time, s.verdict, s.status").
		Where("s.user_id = ? and s.problem_id = ?", params.PerformedByID, params.QuestionId).
		Scan(&res.Response.UserSubmissions).Error;
		
		err != nil {
		return ShowProblemResponse{Status: 400}, err
	}

	return res, nil
}
