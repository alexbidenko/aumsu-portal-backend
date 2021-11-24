package models

import (
	"aumsu.portal.backend/dif"
	"aumsu.portal.backend/entities"
)

type CommentModule struct {
}

func (commentModule CommentModule) Create(model *entities.Comment) {
	dif.DB.Model(&entities.Comment{}).Create(model)
}

func (commentModule CommentModule) Delete(id string) {
	dif.DB.Delete(&entities.Comment{}, id)
}

func (commentModule CommentModule) Update(id string, model *entities.Comment) {
	dif.DB.Model(&entities.Comment{}).Where("id = ?", id).Updates(model)
}
