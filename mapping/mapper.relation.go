package mapping

import (
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"github.com/psyark/notionapi"
)

var _ Mapper = &RelationMapper{}

type RelationMapper struct{}

func (m *RelationMapper) RecordToObject(field reflect.StructField, value reflect.Value, pv *notionapi.PropertyValue) error {
	if _, ok := value.Interface().([]uuid.UUID); ok {
		uuids := []uuid.UUID{}
		for _, rel := range pv.Relation {
			uuids = append(uuids, rel.ID)
		}
		value.Set(reflect.ValueOf(uuids))
		return nil
	} else {
		return fmt.Errorf("unsupported type: %v", field.Type)
	}
}

func (m *RelationMapper) GetDelta(field reflect.StructField, value reflect.Value, pv *notionapi.PropertyValue) (*notionapi.PropertyValue, error) {
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
		return nil, fmt.Errorf("unsupported type: %v", field.Type)
	}
}
