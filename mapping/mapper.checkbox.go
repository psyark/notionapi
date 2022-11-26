package mapping

import (
	"fmt"
	"reflect"

	"github.com/psyark/notionapi"
)

var _ mapper = &checkboxMapper{}

type checkboxMapper struct{ propertyMapper }

func (m *checkboxMapper) decodePage(value reflect.Value, page notionapi.Page) error {
	if _, ok := value.Interface().(bool); ok {
		value.SetBool(m.getPropValue(page).Checkbox)
		return nil
	} else {
		return fmt.Errorf("unsupported type: %v", reflect.TypeOf(value.Interface()))
	}
}

func (m *checkboxMapper) createPageFrom(value reflect.Value, dst *notionapi.CreatePageOptions) error {
	return nil
}

func (m *checkboxMapper) updatePageFrom(value reflect.Value, page notionapi.Page, dst *notionapi.UpdatePageOptions) error {
	return nil
}
