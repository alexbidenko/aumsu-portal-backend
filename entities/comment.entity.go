package entities

type Comment struct {
	Model
	Content   string  `json:"content" valid:"required,type(string),length(1|1024)" gorm:"size:1024"`
	MessageId int     `json:"message_id"`
	UserId    int     `json:"user_id"`
	User      Student `json:"user" gorm:"foreignKey:UserId"`
}
