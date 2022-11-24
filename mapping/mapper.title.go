package mapping

import (
	"fmt"
	"reflect"

	"github.com/psyark/notionapi"
)

var _ Mapper = &TitleMapper{}

type TitleMapper struct{}

func (m *TitleMapper) RecordToObject(field reflect.StructField, value reflect.Value, pv *notionapi.PropertyValue) error {
	if field.Type.Kind() == reflect.String {
		value.SetString(pv.Title.PlainText())
		return nil
	} else if _, ok := value.Interface().(notionapi.RichTextArray); ok {
		value.Set(reflect.ValueOf(pv.Title))
		return nil
	} else {
		return fmt.Errorf("unsupported type: %v", field.Type)
	}
}

func (m *TitleMapper) GetDelta(field reflect.StructField, value reflect.Value, pv *notionapi.PropertyValue) (*notionapi.PropertyValue, error) {
	if field.Type.Kind() == reflect.String {
		if pv == nil || value.String() != pv.Title.PlainText() {
			return &notionapi.PropertyValue{Type: "title", Title: notionapi.RichTextArray{{Type: "text", Text: &notionapi.Text{Content: value.String()}}}}, nil
		}
		return nil, nil
	} else if rta, ok := value.Interface().(notionapi.RichTextArray); ok {
		if pv == nil {
			return &notionapi.PropertyValue{Type: "title", Title: rta}, nil
		}
		if equal, err := compareInJSON(rta, pv.Title); err != nil {
			return nil, err
		} else if !equal {
			return &notionapi.PropertyValue{Type: "title", Title: rta}, nil
		}
		return nil, nil
	} else {
		return nil, fmt.Errorf("unsupported type: %v", field.Type)
	}
}
