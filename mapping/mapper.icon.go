package mapping

import (
	"fmt"
	"reflect"

	"github.com/psyark/notionapi"
)

var _ mapper = &iconMapper{}

type iconMapper struct{}

func (m *iconMapper) decodePage(value reflect.Value, page notionapi.Page) error {
	if value.Type() == reflect.TypeOf((*notionapi.FileOrEmoji)(nil)).Elem() {
		if page.Icon == nil {
			// value.Set(reflect.ValueOf((nil))) // TODO
		} else {
			value.Set(reflect.ValueOf(notionapi.FileOrEmoji(page.Icon)))
		}
		return nil
	} else {
		return fmt.Errorf("unsupported type: %v", reflect.TypeOf(value.Interface()))
	}
}

func (m *iconMapper) createPageFrom(value reflect.Value, dst *notionapi.CreatePageOptions) error {
	delta, err := m.getDelta(value, nil)
	if err != nil {
		return err
	}
	if delta != nil {
		dst.Icon = delta
	}
	return nil
}

func (m *iconMapper) updatePageFrom(value reflect.Value, page notionapi.Page, dst *notionapi.UpdatePageOptions) error {
	delta, err := m.getDelta(value, page.Icon)
	if err != nil {
		return err
	}
	if delta != nil {
		dst.Icon = delta
	}
	return nil
}

func (m *iconMapper) getDelta(value reflect.Value, icon notionapi.FileOrEmoji) (notionapi.FileOrEmoji, error) {
	if value.Type() == reflect.TypeOf((*notionapi.FileOrEmoji)(nil)).Elem() {
		equal, err := compareInJSON(value.Interface(), icon)
		if err != nil {
			return nil, err
		}
		if !equal {
			if foe, ok := value.Interface().(notionapi.FileOrEmoji); ok {
				return foe, nil
			} else {
				return &notionapi.NoIcon{}, nil
			}
		}
		return nil, nil
	} else {
		return nil, fmt.Errorf("unsupported type: %v", reflect.TypeOf(value.Interface()))
	}
}
