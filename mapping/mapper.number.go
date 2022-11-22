package mapping

import (
	"fmt"
	"reflect"

	"github.com/psyark/notionapi"
)

var _ Mapper = &NumberMapper{}

type NumberMapper struct{}

func (m *NumberMapper) RecordToObject(field reflect.StructField, value reflect.Value, pv *notionapi.PropertyValue) error {
	if field.Type.Kind() == reflect.Float64 {
		if pv.Number != nil {
			value.SetFloat(*pv.Number)
		} else {
			value.SetFloat(0)
		}
		return nil
	} else {
		return fmt.Errorf("unsupported type: %v", field.Type)
	}
}
