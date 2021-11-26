package entities

import (
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	Id int `json:"id"`
	Login string `json:"login" valid:"required,type(string),length(5|255)" gorm:"size:255"`
	Password string `json:"password" valid:"required,type(string),length(8|255)" gorm:"size:255"`
	Token string `json:"token" valid:"required,type(string),length(5|255)" gorm:"size:255"`
	FirstName string `json:"firstName" valid:"required,type(string),length(2|255)" gorm:"size:255"`
	LastName string `json:"lastName" valid:"required,type(string),length(2|255)" gorm:"size:255"`
	Avatar string `json:"avatar" valid:"type(string),length(1|255)" gorm:"size:255"`
	Status string `json:"status" valid:"type(string),length(1|255)" gorm:"size:255"`
	Patronymic string `json:"patronymic" valid:"type(string),length(2|255),optional" gorm:"size:255"`
}

func (student *Student) TableName() string {
	return "students"
}

type UserClaims struct {
	Authorized bool
	Student
	jwt.StandardClaims
}

func (student *Student) GenerateJWT() error {
	var signingKey = []byte("aumsu-portal-backend")
	claims := UserClaims{
		Authorized: true,
		Student: *student,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(signingKey)
	student.Token = tokenString

	return err
}
