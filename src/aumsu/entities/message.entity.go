package entities

type Message struct {
	Id int `json:"id"`
	Title string `json:"title" valid:"required,type(string)" gorm:"size:255"`
	Description string `json:"description" valid:"required,type(string)"`
	Image string `json:"image" valid:"required,type(string)" gorm:"size:255"`
	From int `json:"from"`
}

func (faculty *Message) TableName() string {
	return "messages"
}
