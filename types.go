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

type FileOrEmoji interface {
	fileOrEmoji()
}

func newFileOrEmoji(data []byte, childName string) FileOrEmoji {
	t := map[string]interface{}{}
	if err := json.Unmarshal(data, &t); err != nil {
		panic(err)
	}

	if child, ok := t[childName].(map[string]interface{}); ok {
		switch child["type"] {
		case "file":
			return &File{}
		case "emoji":
			return &Emoji{}
		}
	}

	return nil
}

type PropertyItemOrPropertyItemPagination interface {
	propertyItemOrPropertyItemPagination()
}

type PropertyItemOrPropertyItemPaginationUnmarshaller struct {
	value PropertyItemOrPropertyItemPagination
}

func (m *PropertyItemOrPropertyItemPaginationUnmarshaller) UnmarshalJSON(data []byte) error {
	check := struct {
		Object string `json:"object"`
	}{}
	if err := json.Unmarshal(data, &check); err != nil {
		return err
	}

	switch check.Object {
	case "property_item":
		m.value = &PropertyItem{}
	case "list":
		m.value = &PropertyItemPagination{}
	default:
		return fmt.Errorf("unsupported object: %v", check.Object)
	}

	return json.Unmarshal(data, m.value)
}
