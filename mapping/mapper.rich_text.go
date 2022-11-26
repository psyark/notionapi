package mapping

import (
	"fmt"
	"reflect"

	"github.com/psyark/notionapi"
)

var _ mapper = &richTextMapper{}

type richTextMapper struct{ propertyMapper }

func (m *richTextMapper) decodePage(value reflect.Value, page notionapi.Page) error {
	pv := m.getPropValue(page)

	if _, ok := value.Interface().(string); ok {
		value.SetString(pv.RichText.PlainText())
		return nil
	} else if _, ok := value.Interface().(notionapi.RichTextArray); ok {
		value.Set(reflect.ValueOf(pv.RichText))
		return nil
	} else {
		return fmt.Errorf("unsupported type: %v", reflect.TypeOf(value.Interface()))
	}
}

func (m *richTextMapper) createPageFrom(value reflect.Value, dst *notionapi.CreatePageOptions) error {
	delta, err := m.getDelta(value, nil)
	if err != nil {
		return err
	}
	if delta != nil {
		dst.Properties[m.propID] = *delta
	}
	return nil
}

func (m *richTextMapper) updatePageFrom(value reflect.Value, page notionapi.Page, dst *notionapi.UpdatePageOptions) error {
	delta, err := m.getDelta(value, m.getPropValue(page))
	if err != nil {
		return err
	}
	if delta != nil {
		dst.Properties[m.propID] = *delta
	}
	return nil
}

func (m *richTextMapper) getDelta(value reflect.Value, pv *notionapi.PropertyValue) (*notionapi.PropertyValue, error) {
	if _, ok := value.Interface().(string); ok {
		if pv == nil || value.String() != pv.RichText.PlainText() {
			return &notionapi.PropertyValue{Type: "rich_text", RichText: notionapi.RichTextArray{{Type: "text", Text: &notionapi.Text{Content: value.String()}}}}, nil
		}
		return nil, nil
	} else if rta, ok := value.Interface().(notionapi.RichTextArray); ok {
		if pv == nil {
			return &notionapi.PropertyValue{Type: "rich_text", RichText: rta}, nil
		}
		if equal, err := compareInJSON(rta, pv.RichText); err != nil {
			return nil, err
		} else if !equal {
			return &notionapi.PropertyValue{Type: "rich_text", RichText: rta}, nil
		}
		return nil, nil
	} else {
		return nil, fmt.Errorf("unsupported type: %v", reflect.TypeOf(value.Interface()))
	}
}
