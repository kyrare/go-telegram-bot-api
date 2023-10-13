package bot

import "regexp"

type Message struct {
	MessageID int             `json:"message_id"`
	From      User            `json:"from,omitempty"`
	Chat      Chat            `json:"chat,omitempty"`
	Date      int             `json:"date,omitempty"` // timestamp
	Text      string          `json:"text,omitempty"`
	Entities  []MessageEntity `json:"entities,omitempty"`
}

func (m Message) isBotCommand() bool {
	if len(m.Entities) == 0 {
		return false
	}

	return m.Entities[0].Type == "bot_command"
}

func (m Message) BotCommandArgument() string {
	re, _ := regexp.Compile(`^\/[\w\d_]+ (.+)$`)

	subMatches := re.FindStringSubmatch(m.Text)

	if len(subMatches) <= 1 {
		return ""
	}

	return subMatches[1]
}
