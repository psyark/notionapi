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
