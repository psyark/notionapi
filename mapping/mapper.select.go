package mapping

import (
	"fmt"
	"reflect"

	"github.com/psyark/notionapi"
)

var _ mapper = &selectMapper{}

type selectMapper struct{ propertyMapper }

func (m *selectMapper) decodePage(value reflect.Value, page notionapi.Page) error {
	if value.Kind() == reflect.String {
		propValue := m.getPropValue(page)
		if propValue.Select == nil {
			value.SetString("")
		} else {
			value.SetString(propValue.Select.Name)
		}
		return nil
	} else {
		return fmt.Errorf("unsupported type: %v", reflect.TypeOf(value.Interface()))
	}
}

func (m *selectMapper) createPageFrom(value reflect.Value, dst *notionapi.CreatePageOptions) error {
	return nil
}

func (m *selectMapper) updatePageFrom(value reflect.Value, page notionapi.Page, dst *notionapi.UpdatePageOptions) error {
	return nil
}
