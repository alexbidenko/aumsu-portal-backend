package entities

type StudyGroup struct {
	Model
	Name string `json:"name"`
}

func (sg *StudyGroup) TableName() string {
	return "study_groups"
}

type StudyGroupSchedule struct {
	Model
	StudyGroupId uint   `json:"study_group_id"`
	Content      string `json:"content" gorm:"type:text"`
}

func (sgs *StudyGroupSchedule) TableName() string {
	return "study_group_schedules"
}
