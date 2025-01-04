package apis

import (
	"fmt"

	"github.com/lib/pq"
	"github.com/rudrakshsattabhayya/oj_backend_go/config"
	"gorm.io/gorm"
)

type ListProblemsParams struct {
}
type ProblemBasicDetails struct {
	ID                  string         `json:"id"`
	Title               string         `json:"title"`
	Difficulty          int            `json:"difficulty"`
	AcceptedSubmissions int         `json:"acceptedSubmissions"`
	TotalSubmissions    int         `json:"totalSubmissions"`
	Tags                pq.StringArray `gorm:"type:text[]" json:"tags"`
}

type ListProblemsResponse struct {
	Problems []ProblemBasicDetails `json:"problems"`
	Status   string                `json:"status"`
}

func ListProblems(params ListProblemsParams) (ListProblemsResponse, error) {
	db := config.GetDB()

	res, err := GetListProblemsResponse(db)
	if err != nil {
		return ListProblemsResponse{}, err
	}

	return res, nil
}

func GetListProblemsResponse(db *gorm.DB) (ListProblemsResponse, error) {
	res := ListProblemsResponse{}
	problems := []ProblemBasicDetails{}

	query := `
		SELECT p.id, p.title, p.difficulty, p.accepted_submissions, p.total_submissions, ARRAY_AGG(COALESCE(t.name, '')) AS tags
		FROM problems p
		LEFT JOIN problem_tags pt ON pt.problem_id = p.id
		LEFT JOIN tags t ON t.id = pt.tag_id
		GROUP BY p.id, p.title, p.difficulty, p.accepted_submissions, p.total_submissions;
	`
	if err := db.Raw(query).Scan(&problems).Error; err != nil {
		return ListProblemsResponse{}, fmt.Errorf("error executing query: %v", err)
	}

	res.Problems = problems
	res.Status = "success"

	return res, nil
}
