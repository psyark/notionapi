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

func (m *NumberMapper) GetDelta(field reflect.StructField, value reflect.Value, pv *notionapi.PropertyValue) (*notionapi.PropertyValue, error) {
	if field.Type.Kind() == reflect.Float64 {
		if pv.Number == nil || *pv.Number != value.Float() {
			v := value.Float()
			return &notionapi.PropertyValue{Type: "number", Number: &v}, nil
		}
		return nil, nil
	} else {
		return nil, fmt.Errorf("unsupported type: %v", field.Type)
	}
}
