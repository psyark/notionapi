package mapping

import (
	"fmt"
	"reflect"

	"github.com/psyark/notionapi"
)

var _ mapper = &richTextArrayMapper{}

type richTextArrayMapper struct {
	getRTA func(page notionapi.Page) notionapi.RichTextArray
	setRTA func(rta notionapi.RichTextArray, dst notionapi.PropertyValueMap)
}

func (m *richTextArrayMapper) decodePage(value reflect.Value, page notionapi.Page) error {
	rta := m.getRTA(page)

	if _, ok := value.Interface().(string); ok {
		value.SetString(rta.PlainText())
		return nil
	} else if _, ok := value.Interface().(notionapi.RichTextArray); ok {
		value.Set(reflect.ValueOf(rta))
		return nil
	} else {
		return fmt.Errorf("unsupported type: %v", reflect.TypeOf(value.Interface()))
	}
}

func (m *richTextArrayMapper) createPageFrom(value reflect.Value, dst *notionapi.CreatePageOptions) error {
	delta, err := m.getDelta(value, nil)
	if err != nil {
		return err
	}
	if delta != nil {
		m.setRTA(*delta, dst.Properties)
	}
	return nil
}

func (m *richTextArrayMapper) updatePageFrom(value reflect.Value, page notionapi.Page, dst *notionapi.UpdatePageOptions) error {
	rta := m.getRTA(page)
	delta, err := m.getDelta(value, &rta)
	if err != nil {
		return err
	}
	if delta != nil {
		m.setRTA(*delta, dst.Properties)
	}
	return nil
}

func (m *richTextArrayMapper) getDelta(value reflect.Value, src *notionapi.RichTextArray) (*notionapi.RichTextArray, error) {
	if _, ok := value.Interface().(string); ok {
		if src == nil || value.String() != src.PlainText() {
			return &notionapi.RichTextArray{{Type: "text", Text: &notionapi.Text{Content: value.String()}}}, nil
		}
		return nil, nil
	} else if rta, ok := value.Interface().(notionapi.RichTextArray); ok {
		if src == nil {
			return &rta, nil
		}
		if equal, err := compareInJSON(rta, src); err != nil {
			return nil, err
		} else if !equal {
			return &rta, nil
		}
		return nil, nil
	} else {
		return nil, fmt.Errorf("unsupported type: %v", reflect.TypeOf(value.Interface()))
	}
}
