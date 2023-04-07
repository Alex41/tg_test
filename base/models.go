package base

import (
	"time"
)

type Message struct {
	ID        int64 `gorm:"primaryKey"`
	ChatID    int64 `gorm:"primaryKey"`
	Text      string
	FileID    *string
	Creator   MessageCreator `gorm:"type:message_creator"`
	CreatedAt time.Time
}

type MessageCreator string

const (
	MessageCreatorOperator MessageCreator = `operator`
	MessageCreatorUser     MessageCreator = `user`
)
