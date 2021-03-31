package models

import (
	"aumsu/dif"
	"aumsu/entities"
)

type StudentModel struct {
}

func (studentModel StudentModel) Authorization(login string, password string) (entities.Student, error) {
	var student entities.Student
	err := dif.DB.Model(&entities.Student{}).Where(map[string]interface{}{
		"login": login,
		"password": password,
	}).Find(&student).Error

	if err != nil {
		return student, err
	}

	return student, nil
}

func (studentModel StudentModel) GetByToken(token string) (entities.Student, error) {
	var student entities.Student
	err := dif.DB.Model(&entities.Student{}).Where(map[string]interface{}{
		"token": token,
	}).Find(&student).Error

	if err != nil {
		return student, err
	}

	return student, nil
}
