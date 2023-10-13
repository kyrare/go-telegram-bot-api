package bot

import "encoding/json"

type ApiResponse struct {
	Ok        bool            `json:"ok"`
	Result    json.RawMessage `json:"result,omitempty"`
	ErrorCode int             `json:"error_code,omitempty"`
}

type RequestParams map[string]string

type User struct {
	Id        int    `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	UserName  string `json:"username"`
}

type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message,omitempty"`
}

type Chat struct {
	Id int `json:"id"`
}

type MessageEntity struct {
	Offset int    `json:"offset,omitempty"`
	Length int    `json:"length,omitempty"`
	Type   string `json:"type,omitempty"` // bot_command
}
