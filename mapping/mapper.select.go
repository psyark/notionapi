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
	delta, err := m.getDelta(value, nil)
	if err != nil {
		return err
	}
	if delta != nil {
		dst.Properties[m.propID] = *delta
	}
	return nil
}

func (m *selectMapper) updatePageFrom(value reflect.Value, page notionapi.Page, dst *notionapi.UpdatePageOptions) error {
	delta, err := m.getDelta(value, m.getPropValue(page))
	if err != nil {
		return err
	}
	if delta != nil {
		dst.Properties[m.propID] = *delta
	}
	return nil
}

func (m *selectMapper) getDelta(value reflect.Value, pv *notionapi.PropertyValue) (*notionapi.PropertyValue, error) {
	if value.Kind() == reflect.String {
		if value.String() == "" { // あれば消す
			if pv != nil && pv.Select != nil {
				return &notionapi.PropertyValue{Type: "select", Select: nil}, nil
			}
		} else { // 無ければ作る
			if pv == nil || pv.Select == nil || pv.Select.Name != value.String() {
				return &notionapi.PropertyValue{Type: "select", Select: &notionapi.SelectOption{Name: value.String()}}, nil
			}
		}
		return nil, nil
	} else {
		return nil, fmt.Errorf("unsupported type: %v", reflect.TypeOf(value.Interface()))
	}
}
