package dif

import "aumsu.portal.backend/entities"

func Migrate() {
	DB.AutoMigrate(&entities.Student{}, &entities.Message{}, &entities.Comment{})
}
