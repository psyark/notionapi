package mapping

import (
	"fmt"
	"reflect"

	"github.com/psyark/notionapi"
)

type Mapper interface {
	RecordToObject(field reflect.StructField, value reflect.Value, pv *notionapi.PropertyValue) error
}

func getMapper(propType string) (Mapper, error) {
	switch propType {
	case "title":
		return &TitleMapper{}, nil
	case "rich_text":
		return &RichTextMapper{}, nil
	}

	return nil, fmt.Errorf("unsupported property type: %v", propType)
}
