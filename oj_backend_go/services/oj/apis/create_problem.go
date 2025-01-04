package apis

import (
	"fmt"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/rudrakshsattabhayya/oj_backend_go/config"
	"github.com/rudrakshsattabhayya/oj_backend_go/helpers"
	"github.com/rudrakshsattabhayya/oj_backend_go/models/oj"
	"gorm.io/gorm"
)

type CreateProblemParams struct {
	PerformedByID    string                `form:"performed_by_id" binding:"required"`
	Title            string                `form:"title" binding:"required"`
	Difficulty       int                   `form:"difficulty" binding:"required"`
	Tags             string                `form:"tags" binding:"required"`
	ProblemStatement *multipart.FileHeader `form:"problemStatement" binding:"required"`
	HiddenTestCases  *multipart.FileHeader `form:"hiddenTestCases" binding:"required"`
	CorrectSolution  *multipart.FileHeader `form:"correctSolution" binding:"required"`
	VisibleOutputs   *multipart.FileHeader `form:"visibleOutputs" binding:"required"`
	VisibleTestCases *multipart.FileHeader `form:"visibleTestCases" binding:"required"`
}

type CreateProblemResponse struct {
	QuestionId string `json:"QuestionId"`
	Status     string `json:"status"`
}

type CreateProblemUrls struct {
	ProblemStatement string `json:"problemStatement"`
	HiddenTestCases  string `json:"hiddenTestCases"`
	CorrectSolution  string `json:"correctSolution"`
	VisibleOutputs   string `json:"visibleOutputs"`
	VisibleTestCases string `json:"visibleTestCases"`
}

func CreateProblem(params CreateProblemParams) (CreateProblemResponse, error) {
	db := config.GetDB()
	tx := db.Begin()
	res, urls, err := CreateProblemTransCode(params, tx)

	fmt.Println(urls)

	if err != nil {
		//Delete all the uploaded files
		tx.Rollback()
	} else {
		tx.Commit()
	}

	return res, err
}

func CreateProblemTransCode(params CreateProblemParams, tx *gorm.DB) (CreateProblemResponse, CreateProblemUrls, error) {
	res := CreateProblemResponse{
		Status: "400",
	}

	urls, err := UploadFiles(params)
	if err != nil {
		return res, urls, err
	}

	problemID, err := CreateNewProblemInTable(params, urls, tx)
	if err != nil {
		return res, urls, err
	}

	resCreateTags, err := FindOrInitializeTagsTransCode(params.Tags, tx)
	if err != nil {
		return res, urls, err
	}

	err = CreateProblemTags(resCreateTags.Tags, problemID, tx)
	if err != nil {
		return res, urls, err
	}

	// evaluate and save test cases


	res.QuestionId = problemID.String()
	res.Status = "200"

	return res, urls, err
}

func UploadFiles(params CreateProblemParams) (CreateProblemUrls, error) {
	var urls CreateProblemUrls

	funcs := []helpers.FuncWithArgs{
		{ID: "problemStatement", Function: helpers.UploadFile, Args: []interface{}{params.ProblemStatement, "problemStatements"}},
		{ID: "hiddenTestCases", Function: helpers.UploadFile, Args: []interface{}{params.HiddenTestCases, "testCases"}},
		{ID: "correctSolution", Function: helpers.UploadFile, Args: []interface{}{params.CorrectSolution, "solutions"}},
		{ID: "visibleOutputs", Function: helpers.UploadFile, Args: []interface{}{params.VisibleOutputs, "visibleOutputs"}},
		{ID: "visibleTestCases", Function: helpers.UploadFile, Args: []interface{}{params.VisibleTestCases, "visibleTestCases"}},
	}

	results := helpers.RunConcurrently(funcs)

	urls.ProblemStatement = results["problemStatement"].Res.(string)
	urls.HiddenTestCases = results["hiddenTestCases"].Res.(string)
	urls.CorrectSolution = results["correctSolution"].Res.(string)
	urls.VisibleOutputs = results["visibleOutputs"].Res.(string)
	urls.VisibleTestCases = results["visibleTestCases"].Res.(string)

	if err := results["problemStatement"].Err; err != nil {
		return urls, err
	}
	if err := results["hiddenTestCases"].Err; err != nil {
		return urls, err
	}
	if err := results["correctSolution"].Err; err != nil {
		return urls, err
	}
	if err := results["visibleOutputs"].Err; err != nil {
		return urls, err
	}
	if err := results["visibleTestCases"].Err; err != nil {
		return urls, err
	}

	return urls, nil
}

func CreateNewProblemInTable(params CreateProblemParams, urls CreateProblemUrls, tx *gorm.DB) (uuid.UUID, error) {
	problem := oj.Problem{
		Title:            params.Title,
		ProblemStatement: urls.ProblemStatement,
		VisibleTestCases: urls.VisibleTestCases,
		VisibleOutputs:   urls.VisibleOutputs,
		HiddenTestCases:  urls.HiddenTestCases,
		CorrectSolution:  urls.CorrectSolution,
		Difficulty:       params.Difficulty,
	}

	if err := tx.Create(&problem).Error; err != nil {
		return problem.ID, err
	}

	return problem.ID, nil
}

func CreateProblemTags(tags []oj.Tag, problemId uuid.UUID, tx *gorm.DB) error {
	var problemTags []oj.ProblemTag
	for _, Tag := range tags {
		problemTags = append(problemTags, oj.ProblemTag{
			ProblemID: problemId,
			TagID:     Tag.ID,
		})
	}

	if err := tx.Create(&problemTags).Error; err != nil {
		return err
	}

	return nil
}
