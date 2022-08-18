package notionapi

import (
	"encoding/json"
)

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

type PageOrDatabase struct{}

type PropertyItemMarshaler struct {
	PropertyItem
}

func (m *PropertyItemMarshaler) UnmarshalJSON(data []byte) error {
	typeCheck := struct {
		Type string `json:"type"`
	}{}
	if err := json.Unmarshal(data, &typeCheck); err != nil {
		return err
	}

	m.PropertyItem = createPropertyItem(typeCheck.Type)
	return json.Unmarshal(data, m.PropertyItem)
}

func (m *PropertyItemMarshaler) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.PropertyItem)
}
