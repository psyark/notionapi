package mapping

import (
	"fmt"
	"reflect"

	"github.com/psyark/notionapi"
)

var _ mapper = &titleMapper{}

type titleMapper struct{ propertyMapper }

func (m *titleMapper) decodePage(value reflect.Value, page notionapi.Page) error {
	pv := m.getPropValue(page)

	if _, ok := value.Interface().(string); ok {
		value.SetString(pv.Title.PlainText())
		return nil
	} else if _, ok := value.Interface().(notionapi.RichTextArray); ok {
		value.Set(reflect.ValueOf(pv.Title))
		return nil
	} else {
		return fmt.Errorf("unsupported type: %v", reflect.TypeOf(value.Interface()))
	}
}

func (m *titleMapper) createPageFrom(value reflect.Value, dst *notionapi.CreatePageOptions) error {
	delta, err := m.getDelta(value, nil)
	if err != nil {
		return err
	}
	if delta != nil {
		dst.Properties[m.propID] = *delta
	}
	return nil
}

func (m *titleMapper) updatePageFrom(value reflect.Value, page notionapi.Page, dst *notionapi.UpdatePageOptions) error {
	delta, err := m.getDelta(value, m.getPropValue(page))
	if err != nil {
		return err
	}
	if delta != nil {
		dst.Properties[m.propID] = *delta
	}
	return nil
}

func (m *titleMapper) getDelta(value reflect.Value, pv *notionapi.PropertyValue) (*notionapi.PropertyValue, error) {
	if _, ok := value.Interface().(string); ok {
		if pv == nil || value.String() != pv.Title.PlainText() {
			return &notionapi.PropertyValue{Type: "title", Title: notionapi.RichTextArray{{Type: "text", Text: &notionapi.Text{Content: value.String()}}}}, nil
		}
		return nil, nil
	} else if rta, ok := value.Interface().(notionapi.RichTextArray); ok {
		if pv == nil {
			return &notionapi.PropertyValue{Type: "title", Title: rta}, nil
		}
		if equal, err := compareInJSON(rta, pv.Title); err != nil {
			return nil, err
		} else if !equal {
			return &notionapi.PropertyValue{Type: "title", Title: rta}, nil
		}
		return nil, nil
	} else {
		return nil, fmt.Errorf("unsupported type: %v", reflect.TypeOf(value.Interface()))
	}
}
