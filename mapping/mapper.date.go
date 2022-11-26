package mapping

import (
	"fmt"
	"reflect"
	"time"

	"github.com/psyark/notionapi"
)

var _ mapper = &dateMapper{}

type dateMapper struct{ propertyMapper }

func (m *dateMapper) decodePage(value reflect.Value, page notionapi.Page) error {
	pv := m.getPropValue(page)

	if _, ok := value.Interface().(time.Time); ok {
		var t time.Time
		if pv.Date != nil {
			t_, err := m.parseNotionTime(pv.Date.Start)
			if err != nil {
				return err
			}
			t = t_
		}

		value.Set(reflect.ValueOf(t))
		return nil
	} else if _, ok := value.Interface().(*notionapi.DateValue); ok {
		if pv.Date != nil {
			d := *pv.Date
			value.Set(reflect.ValueOf(&d))
		} else {
			value.Set(reflect.ValueOf((*notionapi.DateValue)(nil)))
		}
		return nil
	} else {
		return fmt.Errorf("unsupported type: %v", reflect.TypeOf(value.Interface()))
	}
}

func (m *dateMapper) createPageFrom(value reflect.Value, dst *notionapi.CreatePageOptions) error {
	delta, err := m.getDelta(value, nil)
	if err != nil {
		return err
	}
	if delta != nil {
		dst.Properties[m.propID] = *delta
	}
	return nil
}

func (m *dateMapper) updatePageFrom(value reflect.Value, page notionapi.Page, dst *notionapi.UpdatePageOptions) error {
	delta, err := m.getDelta(value, m.getPropValue(page))
	if err != nil {
		return err
	}
	if delta != nil {
		dst.Properties[m.propID] = *delta
	}
	return nil
}

func (m *dateMapper) getDelta(value reflect.Value, pv *notionapi.PropertyValue) (*notionapi.PropertyValue, error) {
	if timeValue, ok := value.Interface().(time.Time); ok {
		unmatch := false
		if pv == nil {
			unmatch = true
		} else {
			t, err := m.parseNotionTime(pv.Date.Start)
			if err != nil {
				return nil, err
			}
			unmatch = !t.Equal(timeValue)
		}

		if unmatch {
			return &notionapi.PropertyValue{Type: "date", Date: &notionapi.DateValue{Start: timeValue.Format(time.RFC3339Nano)}}, nil
		}
		return nil, nil
	} else if dateValue, ok := value.Interface().(*notionapi.DateValue); ok {
		unmatch := false
		if pv == nil {
			unmatch = true
		} else {
			if equal, err := compareInJSON(dateValue, pv.Date); err != nil {
				return nil, err
			} else if !equal {
				unmatch = true
			}
		}

		if unmatch {
			return &notionapi.PropertyValue{Type: "date", Date: dateValue}, nil
		}
		return nil, nil
	} else {
		return nil, fmt.Errorf("unsupported type: %v", reflect.TypeOf(value.Interface()))
	}
}

func (m *dateMapper) parseNotionTime(timeStr string) (time.Time, error) {
	format := ""
	switch len(timeStr) {
	case len("2022-08-09T00:00:00.000+09:00"):
		format = time.RFC3339Nano
	case len("2006-01-02"):
		format = "2006-01-02"
	default:
		return time.Time{}, fmt.Errorf("unknown format for: %v", timeStr)
	}

	return time.Parse(format, timeStr)
}
