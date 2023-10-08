package models

import (
	"github.com/jinzhu/gorm"
)

type Comment struct {
	gorm.Model
	Text string `json:"text" validate:"required"`
	IsEdited bool `json:"is_edited"`
    EventId int `json:"event_id"`
}

func NewComment() *Comment {
    return &Comment{
        IsEdited: false,
    }
}
