package mapping

import (
	"fmt"
	"reflect"

	"github.com/psyark/notionapi"
)

var _ mapper = &stringPtrMapper{}

type stringPtrMapper struct {
	getStrPtr func(page notionapi.Page) *string
	setStrPtr func(strPtr *string, dst notionapi.PropertyValueMap)
}

func (m *stringPtrMapper) decodePage(value reflect.Value, page notionapi.Page) error {
	strPtr := m.getStrPtr(page)

	if _, ok := value.Interface().(string); ok {
		if strPtr != nil {
			value.SetString(*strPtr)
		} else {
			value.SetString("")
		}
		return nil
	} else {
		return fmt.Errorf("unsupported type: %v", reflect.TypeOf(value.Interface()))
	}
}

func (m *stringPtrMapper) createPageFrom(value reflect.Value, dst *notionapi.CreatePageOptions) error {
	return nil
}

func (m *stringPtrMapper) updatePageFrom(value reflect.Value, page notionapi.Page, dst *notionapi.UpdatePageOptions) error {
	return nil
}
