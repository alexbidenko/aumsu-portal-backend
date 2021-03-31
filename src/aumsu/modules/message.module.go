package models

import (
	"aumsu/dif"
	"aumsu/entities"
)

type MessageModel struct {
}

func (messageModel MessageModel) GetLast() (entities.Message, error) {
	var message entities.Message
	err := dif.DB.Model(&entities.Message{}).Last(&message).Error

	if err != nil {
		return message, err
	}

	return message, nil
}

func (messageModel MessageModel) Create(model *entities.Message) {
	dif.DB.Model(&entities.Message{}).Create(model)
}

func (messageModel MessageModel) All() []entities.Message {
	var messages []entities.Message
	dif.DB.Model(&entities.Message{}).Limit(40).Find(&messages)
	return messages
}
