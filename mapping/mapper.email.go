package mapping

import (
	"fmt"
	"reflect"

	"github.com/psyark/notionapi"
)

var _ mapper = &emailMapper{}

type emailMapper struct{ propertyMapper }

func (m *emailMapper) decodePage(value reflect.Value, page notionapi.Page) error {
	pv := m.getPropValue(page)

	if _, ok := value.Interface().(string); ok {
		if pv.Email != nil {
			value.SetString(*pv.Email)
		} else {
			value.SetString("")
		}
		return nil
	} else {
		return fmt.Errorf("unsupported type: %v", reflect.TypeOf(value.Interface()))
	}
}

func (m *emailMapper) createPageFrom(value reflect.Value, dst *notionapi.CreatePageOptions) error {
	return nil
}

func (m *emailMapper) updatePageFrom(value reflect.Value, page notionapi.Page, dst *notionapi.UpdatePageOptions) error {
	return nil
}
