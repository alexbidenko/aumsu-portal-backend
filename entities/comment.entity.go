package entities

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Id int `json:"id" gorm:"primary_key"`
	Content string `json:"content" valid:"required,type(string),length(1|1024)" gorm:"size:1024"`
	MessageId int `json:"message_id"`
	UserId int `json:"user_id"`
	User Student `json:"user" gorm:"foreignKey:UserId"`
}
