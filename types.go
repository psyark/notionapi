package notionapi

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type UUIDString = string

type ISO8601String = string

type PageReference struct {
	ID string `json:"id,omitempty"`
}

type ImageFile struct {
	File
	Caption []RichText `json:"caption"`
}

func (f ImageFile) MarshalJSON() ([]byte, error) {
	targets := []interface{}{
		f.File,
		struct {
			Caption []RichText `json:"caption"`
		}{f.Caption},
	}

	dst := map[string]interface{}{}
	for _, target := range targets {
		data, err := json.Marshal(target)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, &dst); err != nil {
			return nil, err
		}
	}

	return json.Marshal(dst)
}

// https://developers.notion.com/reference/errors
type Error struct {
	Object  string `json:"object"`
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return fmt.Sprintf("(%v) %v", e.Code, e.Message)
}

type PageOrDatabase struct{}

func marshalByType(object interface{}, typeValue string) ([]byte, error) {
	data, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}

	m := map[string]interface{}{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	typ := reflect.TypeOf(object)
	for i := 0; i < typ.NumField(); i++ {
		fld := typ.Field(i)
		propName := strings.TrimSuffix(fld.Tag.Get("json"), ",omitempty")
		if fld.Tag.Get("specific") == "type" && typeValue != propName {
			delete(m, propName)
		}
	}

	return json.Marshal(m)
}

type FileOrEmoji struct {
	File  *File
	Emoji *Emoji
}

func (u *FileOrEmoji) UnmarshalJSON(data []byte) error {
	check := struct {
		Type string `json:"type"`
	}{}
	if err := json.Unmarshal(data, &check); err != nil {
		return err
	}

	u.File = nil
	u.Emoji = nil

	switch check.Type {
	case "file":
		u.File = &File{}
		return json.Unmarshal(data, u.File)
	case "emoji":
		u.Emoji = &Emoji{}
		return json.Unmarshal(data, u.Emoji)
	default:
		return fmt.Errorf("unknown type: %v", check.Type)
	}
}

func (u FileOrEmoji) MarshalJSON() ([]byte, error) {
	switch {
	case u.File != nil:
		return json.Marshal(u.File)
	case u.Emoji != nil:
		return json.Marshal(u.Emoji)
	default:
		return nil, fmt.Errorf("at least one must be non-nil: file, emoji")
	}
}

type PropertyItemOrPagination struct {
	Object       string `json:"object"`
	PropertyItem PropertyItem
	Pagination   PropertyItemPagination
}

func (m *PropertyItemOrPagination) UnmarshalJSON(data []byte) error {
	check := struct {
		Object string `json:"object"`
	}{}
	if err := json.Unmarshal(data, &check); err != nil {
		return err
	}

	m.Object = check.Object
	switch m.Object {
	case "property_item":
		return json.Unmarshal(data, &m.PropertyItem)
	case "list":
		return json.Unmarshal(data, &m.Pagination)
	default:
		panic(m.Object)
	}
}

func (m PropertyItemOrPagination) MarshalJSON() ([]byte, error) {
	switch m.Object {
	case "property_item":
		return json.Marshal(m.PropertyItem)
	case "list":
		return json.Marshal(m.Pagination)
	default:
		panic(m.Object)
	}
}

func (p PropertyItemPagination) MarshalJSON() ([]byte, error) {
	itemBytes, err := json.Marshal(p.PropertyItem)
	if err != nil {
		return nil, err
	}

	itemMap := map[string]interface{}{}
	if err := json.Unmarshal(itemBytes, &itemMap); err != nil {
		return nil, err
	}

	delete(itemMap, "object")
	for k := range itemMap {
		if k != "id" && k != "type" && k != "rollup" {
			itemMap[k] = struct{}{} // rollup以外のタイプごとのプロパティは全て {} にする
		}
	}
	itemMap["next_url"] = p.PropertyItem.NextURL

	return json.Marshal(struct {
		Pagination
		Results      []PropertyItem         `json:"results"`
		PropertyItem map[string]interface{} `json:"property_item"`
	}{
		p.Pagination,
		p.Results,
		itemMap,
	})
}
