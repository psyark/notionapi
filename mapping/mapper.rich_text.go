package mapping

import (
	"fmt"
	"reflect"

	"github.com/psyark/notionapi"
)

var _ Mapper = &RichTextMapper{}

type RichTextMapper struct{}

func (m *RichTextMapper) RecordToObject(field reflect.StructField, value reflect.Value, pv *notionapi.PropertyValue) error {
	if field.Type.Kind() == reflect.String {
		value.SetString(pv.RichText.PlainText())
		return nil
	} else if _, ok := value.Interface().(notionapi.RichTextArray); ok {
		value.Set(reflect.ValueOf(pv.RichText))
		return nil
	} else {
		return fmt.Errorf("unsupported type: %v", field.Type)
	}
}

func (m *RichTextMapper) GetDelta(field reflect.StructField, value reflect.Value, pv *notionapi.PropertyValue) (*notionapi.PropertyValue, error) {
	if field.Type.Kind() == reflect.String {
		if value.String() != pv.RichText.PlainText() {
			return &notionapi.PropertyValue{Type: "rich_text", RichText: notionapi.RichTextArray{{Type: "text", Text: &notionapi.Text{Content: value.String()}}}}, nil
		}
		return nil, nil
	} else if rta, ok := value.Interface().(notionapi.RichTextArray); ok {
		if equal, err := compareInJSON(rta, pv.RichText); err != nil {
			return nil, err
		} else if !equal {
			return &notionapi.PropertyValue{Type: "rich_text", RichText: rta}, nil
		}
		return nil, nil
	} else {
		return nil, fmt.Errorf("unsupported type: %v", field.Type)
	}
}
