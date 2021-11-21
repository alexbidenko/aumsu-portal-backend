package dif

import "aumsu/entities"

func Migrate() {
	DB.AutoMigrate(&entities.Student{}, &entities.Message{}, &entities.Comment{})
}
