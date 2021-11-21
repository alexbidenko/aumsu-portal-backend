package models

import (
	"aumsu/dif"
	"aumsu/entities"
)

type CommentModule struct {
}

func (commentModule CommentModule) Create(model *entities.Comment) {
	dif.DB.Model(&entities.Comment{}).Create(model)
}
