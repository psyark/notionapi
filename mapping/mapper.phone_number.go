package mapping

import (
	"fmt"
	"reflect"

	"github.com/psyark/notionapi"
)

var _ Mapper = &PhoneNumberMapper{}

type PhoneNumberMapper struct{}

func (m *PhoneNumberMapper) RecordToObject(field reflect.StructField, value reflect.Value, pv *notionapi.PropertyValue) error {
	if field.Type.Kind() == reflect.String {
		if pv.PhoneNumber != nil {
			value.SetString(*pv.PhoneNumber)
		} else {
			value.SetString("")
		}
		return nil
	} else {
		return fmt.Errorf("unsupported type: %v", field.Type)
	}
}
