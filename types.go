package notionapi

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

var ErrTypeUnset = errors.New("err type unset")

type UUIDString = string

type ISO8601String = string

type PageReference struct {
	ID string `json:"id,omitempty"`
}

// https://developers.notion.com/reference/errors
type Error struct {
	Object  string `json:"object"`
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type PageOrDatabase struct{}

func marshalByType(object interface{}, typeValue string) ([]byte, error) {
	if typeValue == "" {
		return nil, ErrTypeUnset
	}

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
