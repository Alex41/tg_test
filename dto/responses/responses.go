package responses

import "time"

type Message struct {
	ID        int64     `json:"id"`
	ChatID    int64     `json:"chat_id"`
	Text      string    `json:"text"`
	File      *string   `json:"file"`
	Creator   string    `json:"creator"`
	CreatedAt time.Time `json:"created_at"`
}
