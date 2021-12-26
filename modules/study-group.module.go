package models

import (
	"aumsu.portal.backend/dif"
	"aumsu.portal.backend/entities"
	"encoding/json"
	"fmt"
)

type StudyGroupModel struct {
}

func (StudyGroupModel) All() []entities.StudyGroup {
	studyGroups := make([]entities.StudyGroup, 0)
	dif.DB.Model(&entities.StudyGroup{}).Find(&studyGroups)
	return studyGroups
}

func (StudyGroupModel) GetSchedule(studyGroupId string) (map[string]interface{}, error) {
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
