package mapping

import (
	"fmt"
	"reflect"

	"github.com/psyark/notionapi"
)

var _ mapper = &urlMapper{}

type urlMapper struct{ propertyMapper }

func (m *urlMapper) decodePage(value reflect.Value, page notionapi.Page) error {
	pv := m.getPropValue(page)

	if _, ok := value.Interface().(string); ok {
		if pv.URL != nil {
			value.SetString(*pv.URL)
		} else {
			value.SetString("")
		}
		return nil
	} else {
		return fmt.Errorf("unsupported type: %v", reflect.TypeOf(value.Interface()))
	}
}

func (m *urlMapper) createPageFrom(value reflect.Value, dst *notionapi.CreatePageOptions) error {
	return nil
}

func (m *urlMapper) updatePageFrom(value reflect.Value, page notionapi.Page, dst *notionapi.UpdatePageOptions) error {
	return nil
}
