package models

import (
	"aumsu.portal.backend/dif"
	"aumsu.portal.backend/entities"
)

type StudyGroupModel struct {
}

func (StudyGroupModel) All() []entities.StudyGroup {
	studyGroups := make([]entities.StudyGroup, 0)
	dif.DB.Model(&entities.StudyGroup{}).Find(&studyGroups)
	return studyGroups
}

func (StudyGroupModel) GetSchedule(studyGroupId string) (entities.StudyGroupSchedule, error) {
	var studyGroupSchedule entities.StudyGroupSchedule
	err := dif.DB.Model(&entities.StudyGroupSchedule{}).Where("study_group_id", studyGroupId).First(&studyGroupSchedule).Error
	return studyGroupSchedule, err
}
