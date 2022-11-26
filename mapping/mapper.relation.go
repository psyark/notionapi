package mapping

import (
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"github.com/psyark/notionapi"
)

var _ mapper = &relationMapper{}

type relationMapper struct{ propertyMapper }

func (m *relationMapper) decodePage(value reflect.Value, page notionapi.Page) error {
	pv := m.getPropValue(page)

	if _, ok := value.Interface().([]uuid.UUID); ok {
		uuids := []uuid.UUID{}
		for _, rel := range pv.Relation {
			uuids = append(uuids, rel.ID)
		}
		value.Set(reflect.ValueOf(uuids))
		return nil
	} else {
		return fmt.Errorf("unsupported type: %v", reflect.TypeOf(value.Interface()))
	}
}

func (m *relationMapper) createPageFrom(value reflect.Value, dst *notionapi.CreatePageOptions) error {
	delta, err := m.getDelta(value, nil)
	if err != nil {
		return err
	}
	if delta != nil {
		dst.Properties[m.propID] = *delta
	}
	return nil
}

func (m *relationMapper) updatePageFrom(value reflect.Value, page notionapi.Page, dst *notionapi.UpdatePageOptions) error {
	delta, err := m.getDelta(value, m.getPropValue(page))
	if err != nil {
		return err
	}
	if delta != nil {
		dst.Properties[m.propID] = *delta
	}
	return nil
}

func (m *relationMapper) getDelta(value reflect.Value, pv *notionapi.PropertyValue) (*notionapi.PropertyValue, error) {
	if uuids, ok := value.Interface().([]uuid.UUID); ok {
		unmatch := false
		if pv == nil {
			unmatch = true
		} else {
			uuids2 := []uuid.UUID{}
			for _, rel := range pv.Relation {
				uuids2 = append(uuids2, rel.ID)
			}
			if equal, err := compareInJSON(uuids, uuids2); err != nil {
				return nil, err
			} else if !equal {
				unmatch = true
			}
		}

		if unmatch {
			pv2 := &notionapi.PropertyValue{Type: "relation"}
			for _, u := range uuids {
				pv2.Relation = append(pv2.Relation, notionapi.PageReference{ID: u})
			}
			return pv2, nil
		}
		return nil, nil
	} else {
		return nil, fmt.Errorf("unsupported type: %v", reflect.TypeOf(value.Interface()))
	}
}
