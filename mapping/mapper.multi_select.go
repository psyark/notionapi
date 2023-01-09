package mapping

import (
	"fmt"
	"reflect"

	"github.com/psyark/notionapi"
)

var _ mapper = &multiSelectMapper{}

type multiSelectMapper struct{ propertyMapper }

func (m *multiSelectMapper) decodePage(value reflect.Value, page notionapi.Page) error {
	if _, ok := value.Interface().([]string); ok {
		propValue := m.getPropValue(page)
		strs := make([]string, len(propValue.MultiSelect))
		for i, opt := range propValue.MultiSelect {
			strs[i] = opt.Name
		}
		value.Set(reflect.ValueOf(strs))
		return nil
	} else {
		return fmt.Errorf("unsupported type: %v", reflect.TypeOf(value.Interface()))
	}
}

func (m *multiSelectMapper) createPageFrom(value reflect.Value, dst *notionapi.CreatePageOptions) error {
	return nil
}

func (m *multiSelectMapper) updatePageFrom(value reflect.Value, page notionapi.Page, dst *notionapi.UpdatePageOptions) error {
	return nil
}
