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
	Caption []RichText `json:"caption,omitempty"` // image.caption should be an array or undefined, instead was null
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

var (
	_ FileOrEmoji = &File{}
	_ FileOrEmoji = &Emoji{}
)

type FileOrEmoji interface {
	fileOrEmoji()
}

func newFileOrEmoji(data []byte) FileOrEmoji {
	switch string(getChild(data, "type")) {
	case `"file"`:
		return &File{}
	case `"emoji"`:
		return &Emoji{}
	}

	return nil
}

var (
	_ PropertyItemOrPropertyItemPagination = &PropertyItem{}
	_ PropertyItemOrPropertyItemPagination = &PropertyItemPagination{}
)

type PropertyItemOrPropertyItemPagination interface {
	propertyItemOrPropertyItemPagination()
}

type _PropertyItemOrPropertyItemPaginationUnmarshaller struct {
	PropertyItemOrPropertyItemPagination
}

func (m *_PropertyItemOrPropertyItemPaginationUnmarshaller) UnmarshalJSON(data []byte) error {
	switch string(getChild(data, "object")) {
	case `"property_item"`:
		m.PropertyItemOrPropertyItemPagination = &PropertyItem{}
	case `"list"`:
		m.PropertyItemOrPropertyItemPagination = &PropertyItemPagination{}
	default:
		return fmt.Errorf("unsupported object")
	}

	return json.Unmarshal(data, m.PropertyItemOrPropertyItemPagination)
}

func getChild(data []byte, childName string) []byte {
	t := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &t); err != nil {
		panic(err)
	}

	return t[childName]
}
