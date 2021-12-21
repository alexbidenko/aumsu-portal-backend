package models

import (
	"aumsu.portal.backend/dif"
	"aumsu.portal.backend/entities"
)

type StudentModel struct {
}

func (studentModel StudentModel) Authorization(login string) (entities.Student, error) {
	var student entities.Student
	err := dif.DB.Model(&entities.Student{}).Where("login = ?", login).Preload("StudyGroup").First(&student).Error

	if err != nil {
		return student, err
	}

	return student, nil
}

func (studentModel StudentModel) Create(student *entities.Student) error {
	return dif.DB.Model(&entities.Student{}).Create(student).Error
}

func (studentModel StudentModel) CheckUnique(student *entities.Student) bool {
	err := dif.DB.Model(&entities.Student{}).Where("login = ?", student.Login).Find(&student).Error

	return err != nil
}

func (studentModel StudentModel) GetByToken(token string) (entities.Student, error) {
	var student entities.Student
	err := dif.DB.Model(&entities.Student{}).Where(map[string]interface{}{
		"token": token,
	}).Preload("StudyGroup").Find(&student).Error

	if err != nil {
		return student, err
	}

	return student, nil
}

func (studentModel StudentModel) Update(id uint, student *entities.Student) {
	dif.DB.Model(&entities.Student{}).
		Where("id = ?", id).
		Updates(student)
	dif.DB.Model(&entities.Student{}).
		Where("id = ?", id).
		Preload("StudyGroup").
		First(student)
}
