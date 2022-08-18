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

type PropertyItems []PropertyItem

func (items *PropertyItems) UnmarshalJSON(data []byte) error {
	type typeCheck struct {
		Type string `json:"type"`
	}
	tcs := []typeCheck{}
	if err := json.Unmarshal(data, &tcs); err != nil {
		return err
	}

	*items = make([]PropertyItem, len(tcs))
	for i, tc := range tcs {
		(*items)[i] = createPropertyItem(tc.Type)
	}

	type Alias PropertyItems
	if err := json.Unmarshal(data, (*Alias)(items)); err != nil {
		return err
	}

	return nil
}
