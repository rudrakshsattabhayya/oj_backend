package apis

import (
	"fmt"
	"strings"

	"github.com/rudrakshsattabhayya/oj_backend_go/config"
	"github.com/rudrakshsattabhayya/oj_backend_go/helpers"
	"github.com/rudrakshsattabhayya/oj_backend_go/models/oj"
	"gorm.io/gorm"
)

type CreateTagsParams struct {
	Tags string `json:"tags" binding:"required"`
}

type CreateTagsResponse struct {
	Tags []oj.Tag `json:"tags"`
}

func CreateTags(params CreateTagsParams) (CreateTagsResponse, error) {
	db := config.GetDB()
	tx := db.Begin()

	res, err := FindOrInitializeTagsTransCode(params.Tags, tx)

	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	return res, err
}

func FindOrInitializeTagsTransCode(tagsString string, tx *gorm.DB) (CreateTagsResponse, error) {
	tags := strings.Split(tagsString, ", ")
	if len(tags) == 0 {
		return CreateTagsResponse{}, fmt.Errorf("tags cannot be empty")
	}

	tags, existingTagModels, err := GetUncreatedTags(tags, tx)
	if err != nil {
		return CreateTagsResponse{}, err
	}

	res, err := BatchInsertTags(tags, tx)
	if err != nil {
		return CreateTagsResponse{}, err
	}

	res.Tags = append(res.Tags, existingTagModels...)

	return res, nil
}

func GetUncreatedTags(tags []string, tx *gorm.DB) ([]string, []oj.Tag, error) {
	var uncreatedTags, existingTags []string
	var existingTagModels []oj.Tag

	if result := tx.Model(&oj.Tag{}).Where("name IN ?", tags).Find(&existingTagModels); result.Error != nil {
		return nil, existingTagModels, result.Error
	}

	for _, Tag := range existingTagModels {
		existingTags = append(existingTags, Tag.Name)
	}

	uncreatedTags = helpers.Difference(tags, existingTags)

	return uncreatedTags, existingTagModels, nil
}

func BatchInsertTags(tags []string, tx *gorm.DB) (CreateTagsResponse, error) {
	var tagModels []oj.Tag
	for _, tag := range tags {
		tagModels = append(tagModels, oj.Tag{Name: tag})
	}

	if len(tagModels) != 0 {
		if err := tx.Create(&tagModels).Error; err != nil {
			return CreateTagsResponse{}, err
		}
	}

	return CreateTagsResponse{Tags: tagModels}, nil
}
