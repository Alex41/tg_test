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

// I used enum for store creator as string with low memory loss on save

type MessageCreator string

const (
	MessageCreatorOperator MessageCreator = `operator`
	MessageCreatorUser     MessageCreator = `user`
)
