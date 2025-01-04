package apis

import (
	"github.com/rudrakshsattabhayya/oj_backend_go/config"
	"gorm.io/gorm"
)

type ListTagsParams struct {
}

type ListTagsResponse struct {
	Tags []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"tags"`
	Status int `json:"status"`
}

func ListTags(params ListTagsParams) (ListTagsResponse, error) {
	db := config.GetDB()

	res, err := GetListTags(db)
	if err != nil {
		return ListTagsResponse{Status: 400}, err
	}

	return res, nil
}

func GetListTags(db *gorm.DB) (ListTagsResponse, error) {
	res := ListTagsResponse{Status: 200}

	if err := db.Table("tags").Select("id, name").Scan(&res.Tags).Error; err != nil {
		return ListTagsResponse{Status: 400}, err
	}

	return res, nil
}
