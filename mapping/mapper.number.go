package mapping

import (
	"fmt"
	"reflect"

	"github.com/psyark/notionapi"
)

var _ mapper = &numberMapper{}

type numberMapper struct{ propertyMapper }

func (m *numberMapper) decodePage(value reflect.Value, page notionapi.Page) error {
	pv := m.getPropValue(page)

	if _, ok := value.Interface().(float64); ok {
		if pv.Number != nil {
			value.SetFloat(*pv.Number)
		} else {
			value.SetFloat(0)
		}
		return nil
	} else {
		return fmt.Errorf("unsupported type: %v", reflect.TypeOf(value.Interface()))
	}
}

func (m *numberMapper) createPageFrom(value reflect.Value, dst *notionapi.CreatePageOptions) error {
	delta, err := m.getDelta(value, nil)
	if err != nil {
		return err
	}
	if delta != nil {
		dst.Properties[m.propID] = *delta
	}
	return nil
}

func (m *numberMapper) updatePageFrom(value reflect.Value, page notionapi.Page, dst *notionapi.UpdatePageOptions) error {
	delta, err := m.getDelta(value, m.getPropValue(page))
	if err != nil {
		return err
	}
	if delta != nil {
		dst.Properties[m.propID] = *delta
	}
	return nil
}

func (m *numberMapper) getDelta(value reflect.Value, pv *notionapi.PropertyValue) (*notionapi.PropertyValue, error) {
	if _, ok := value.Interface().(float64); ok {
		if pv == nil || pv.Number == nil || *pv.Number != value.Float() {
			v := value.Float()
			return &notionapi.PropertyValue{Type: "number", Number: &v}, nil
		}
		return nil, nil
	} else {
		return nil, fmt.Errorf("unsupported type: %v", reflect.TypeOf(value.Interface()))
	}
}
