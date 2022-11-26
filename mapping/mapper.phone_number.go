package mapping

import (
	"fmt"
	"reflect"

	"github.com/psyark/notionapi"
)

var _ mapper = &phoneNumberMapper{}

type phoneNumberMapper struct{ propertyMapper }

func (m *phoneNumberMapper) decodePage(value reflect.Value, page notionapi.Page) error {
	pv := m.getPropValue(page)

	if _, ok := value.Interface().(string); ok {
		if pv.PhoneNumber != nil {
			value.SetString(*pv.PhoneNumber)
		} else {
			value.SetString("")
		}
		return nil
	} else {
		return fmt.Errorf("unsupported type: %v", reflect.TypeOf(value.Interface()))
	}
}

func (m *phoneNumberMapper) createPageFrom(value reflect.Value, dst *notionapi.CreatePageOptions) error {
	return nil
}

func (m *phoneNumberMapper) updatePageFrom(value reflect.Value, page notionapi.Page, dst *notionapi.UpdatePageOptions) error {
	return nil
}
