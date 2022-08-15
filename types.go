package notionapi

type UUIDString string

type ISO8601String string

type FileOrEmoji struct {
	Type string `json:"type"`
	*File
	*Emoji
}

// https://developers.notion.com/reference/errors
type Error struct {
	Object  string `json:"object"`
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}
