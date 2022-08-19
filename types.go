package notionapi

import (
	"encoding/json"
)

type UUIDString string

type ISO8601String string

type PageReference struct {
	ID string `json:"id,omitempty"`
}

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

type PropertyItemOrPagination struct {
	PropertyItemMarshaler
	PropertyItemPagination
	Object string `json:"object"`
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
		return json.Unmarshal(data, &m.PropertyItemMarshaler)
	case "list":
		return json.Unmarshal(data, &m.PropertyItemPagination)
	default:
		panic(m.Object)
	}
}

func (m PropertyItemOrPagination) MarshalJSON() ([]byte, error) {
	switch m.Object {
	case "property_item":
		return json.Marshal(m.PropertyItemMarshaler)
	case "list":
		return json.Marshal(m.PropertyItemPagination)
	default:
		panic(m.Object)
	}
}

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

	pi, err := createPropertyItem(typeCheck.Type)
	if err != nil {
		return err
	}

	m.PropertyItem = pi
	return json.Unmarshal(data, m.PropertyItem)
}

func (m PropertyItemMarshaler) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.PropertyItem)
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
	itemMap["next_url"] = p.PropertyItem.getCommon().NextURL

	return json.Marshal(struct {
		Pagination
		Results      []PropertyItemMarshaler `json:"results"`
		PropertyItem map[string]interface{}  `json:"property_item"`
	}{
		p.Pagination,
		p.Results,
		itemMap,
	})
}
