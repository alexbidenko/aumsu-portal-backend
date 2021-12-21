package entities

type Message struct {
	Model
	Title       string    `json:"title" valid:"required,type(string)" gorm:"size:255"`
	Description string    `json:"description" valid:"required,type(string)"`
	Image       string    `json:"image" valid:"required,type(string)" gorm:"size:255"`
	From        uint      `json:"from"`
	Comments    []Comment `json:"comments,omitempty" gorm:"foreignKey:MessageId"`
}

func (faculty *Message) TableName() string {
	return "messages"
}
