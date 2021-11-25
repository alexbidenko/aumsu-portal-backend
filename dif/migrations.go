package dif

import (
	"aumsu.portal.backend/entities"
	models "aumsu.portal.backend/modules"
	"golang.org/x/crypto/bcrypt"
)

func Migrate() {
	DB.AutoMigrate(&entities.Student{}, &entities.Message{}, &entities.Comment{})
	
	var student entities.Student
	DB.Model(&entities.Student{}).Where("login = ?", "alexbidenko").First(&student)
	if student.Password == "12345678" {
		var students []entities.Student
		DB.Model(&entities.Student{}).Find(&students)

		for _, item := range students {
			bytes, _ := bcrypt.GenerateFromPassword([]byte(item.Password), 14)
			item.Password = string(bytes)

			var studentModule models.StudentModel
			studentModule.Update(item.Id, &item)
		}
	}
}
