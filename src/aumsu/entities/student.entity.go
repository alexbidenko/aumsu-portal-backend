package entities

type Student struct {
	Id int `json:"id"`
	Login string `json:"login" valid:"required,type(string),length(1|255)" gorm:"size:255"`
	Password string `json:"password" valid:"required,type(string),length(1|255)" gorm:"size:255"`
	Token string `json:"token" valid:"required,type(string),length(1|255)" gorm:"size:255"`
	FirstName string `json:"firstName" valid:"type(string),length(1|255)" gorm:"size:255"`
	LastName string `json:"lastName" valid:"type(string),length(1|255)" gorm:"size:255"`
	Avatar string `json:"avatar" valid:"type(string),length(1|255)" gorm:"size:255"`
	Status string `json:"status" valid:"type(string),length(1|255)" gorm:"size:255"`
	Patronymic string `json:"patronymic" valid:"type(string),length(1|255)" gorm:"size:255"`
}

func (faculty *Student) TableName() string {
	return "students"
}
