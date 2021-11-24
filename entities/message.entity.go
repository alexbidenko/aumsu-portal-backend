package entities

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	Id int `json:"id" gorm:"primary_key"`
	Title string `json:"title" valid:"required,type(string)" gorm:"size:255"`
	Description string `json:"description" valid:"required,type(string)"`
	Image string `json:"image" valid:"required,type(string)" gorm:"size:255"`
	From int `json:"from"`
	Comments []Comment `json:"comments,omitempty" gorm:"foreignKey:MessageId"`
}

func (faculty *Message) TableName() string {
	return "messages"
}
