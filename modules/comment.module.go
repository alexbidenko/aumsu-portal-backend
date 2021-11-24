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
