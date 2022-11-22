package mapping

import (
	"fmt"
	"reflect"
	"time"

	"github.com/psyark/notionapi"
)

var _ Mapper = &DateMapper{}

type DateMapper struct{}

func (m *DateMapper) RecordToObject(field reflect.StructField, value reflect.Value, pv *notionapi.PropertyValue) error {
	if _, ok := value.Interface().(time.Time); ok {
		if pv.Date != nil {
			format := ""
			switch len(pv.Date.Start) {
			case len("2022-08-09T00:00:00.000+09:00"):
				format = time.RFC3339Nano
			case len("2006-01-02"):
				format = "2006-01-02"
			default:
				return fmt.Errorf("unknown format for: %v", pv.Date.Start)
			}
			d, err := time.Parse(format, pv.Date.Start)
			if err != nil {
				return err
			}
			value.Set(reflect.ValueOf(d))
		} else {
			value.Set(reflect.ValueOf(time.Time{}))
		}
		return nil
	} else if _, ok := value.Interface().(*notionapi.DateValue); ok {
		value.Set(reflect.ValueOf(pv.Date))
		return nil
	} else {
		return fmt.Errorf("unsupported type: %v", field.Type)
	}
}

func (m *DateMapper) GetDelta(field reflect.StructField, value reflect.Value, pv *notionapi.PropertyValue) (*notionapi.PropertyValue, error) {
	return nil, nil
}
