package models

import (
	"aumsu.portal.backend/dif"
	"aumsu.portal.backend/entities"
	"encoding/json"
	"errors"
	"fmt"
)

type StudyGroupModel struct {
}

func (StudyGroupModel) All() []entities.StudyGroup {
	studyGroups := make([]entities.StudyGroup, 0)
	dif.DB.Model(&entities.StudyGroup{}).Find(&studyGroups)
	return studyGroups
}

func (StudyGroupModel) GetByGroupId(studyGroupId string) (map[string]interface{}, error) {
	var studyGroupSchedule map[string]interface{}
	err := dif.DB.Model(&entities.StudyGroupSchedule{}).Where("study_group_id", studyGroupId).First(&studyGroupSchedule).Error

	if err == nil {
		var data interface{}
		err = json.Unmarshal([]byte(fmt.Sprintf("%v", studyGroupSchedule["content"])), &data)

		if err == nil {
			studyGroupSchedule["content"] = data
		}
	}

	return studyGroupSchedule, err
}

func (StudyGroupModel) CreateByGroupId(studyGroupId uint, data string) (entities.StudyGroupSchedule, error) {
	var count int64
	dif.DB.Model(&entities.StudyGroupSchedule{}).Where("study_group_id", studyGroupId).Count(&count)
	if count > 0 {
		return entities.StudyGroupSchedule{}, errors.New("schedule already exists")
	}

	studyGroupSchedule := entities.StudyGroupSchedule{
		StudyGroupId: studyGroupId,
		Content:      data,
	}
	err := dif.DB.Model(&entities.StudyGroupSchedule{}).Create(&studyGroupSchedule).Error

	return studyGroupSchedule, err
}

func (StudyGroupModel) UpdateByGroupId(studyGroupId string, data string) error {
	err := dif.DB.Model(&entities.StudyGroupSchedule{}).Where("study_group_id = ?", studyGroupId).Update("content", data).Error

	return err
}
