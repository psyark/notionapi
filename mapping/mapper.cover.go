package mapping

import (
	"fmt"
	"reflect"

	"github.com/psyark/notionapi"
)

var _ mapper = &coverMapper{}

type coverMapper struct{}

func (m *coverMapper) decodePage(value reflect.Value, page notionapi.Page) error {
	if _, ok := value.Interface().(*notionapi.File); ok {
		value.Set(reflect.ValueOf(page.Cover))
		return nil
	} else {
		return fmt.Errorf("unsupported type: %v", reflect.TypeOf(value.Interface()))
	}
}

func (m *coverMapper) createPageFrom(value reflect.Value, dst *notionapi.CreatePageOptions) error {
	delta, err := m.getDelta(value, nil)
	if err != nil {
		return err
	}
	if delta != nil {
		dst.Cover = delta
	}
	return nil
}

func (m *coverMapper) updatePageFrom(value reflect.Value, page notionapi.Page, dst *notionapi.UpdatePageOptions) error {
	delta, err := m.getDelta(value, page.Cover)
	if err != nil {
		return err
	}
	if delta != nil {
		dst.Cover = delta
	}
	return nil
}

func (m *coverMapper) getDelta(value reflect.Value, file *notionapi.File) (*notionapi.File, error) {
	if fileValue, ok := value.Interface().(*notionapi.File); ok {
		equal, err := compareInJSON(fileValue, file)
		if err != nil {
			return nil, err
		}
		if !equal {
			return fileValue, nil
		}
		return nil, nil
	} else {
		return nil, fmt.Errorf("unsupported type: %v", reflect.TypeOf(value.Interface()))
	}
}
