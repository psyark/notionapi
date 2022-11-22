package mapping

import (
	"fmt"
	"reflect"

	"github.com/psyark/notionapi"
)

var _ Mapper = &CheckboxMapper{}

type CheckboxMapper struct{}

func (m *CheckboxMapper) RecordToObject(field reflect.StructField, value reflect.Value, pv *notionapi.PropertyValue) error {
	if field.Type.Kind() == reflect.Bool {
		value.SetBool(pv.Checkbox)
		return nil
	} else {
		return fmt.Errorf("unsupported type: %v", field.Type)
	}
}
