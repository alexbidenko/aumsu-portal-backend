package entities

type Message struct {
	Id int `json:"id"`
	Message string `json:"message" valid:"required,type(string)"`
	From int `json:"from"`
}

func (faculty *Message) TableName() string {
	return "messages"
}
