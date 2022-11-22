package mapping

import (
	"fmt"
	"reflect"

	"github.com/psyark/notionapi"
)

type Mapper interface {
	RecordToObject(field reflect.StructField, value reflect.Value, pv *notionapi.PropertyValue) error
	GetDelta(field reflect.StructField, value reflect.Value, pv *notionapi.PropertyValue) (*notionapi.PropertyValue, error)
}

func getMapper(propType string) (Mapper, error) {
	switch propType {
	case "title":
		return &TitleMapper{}, nil
	case "rich_text":
		return &RichTextMapper{}, nil
	case "email":
		return &EmailMapper{}, nil
	case "url":
		return &URLMapper{}, nil
	case "phone_number":
		return &PhoneNumberMapper{}, nil
	case "number":
		return &NumberMapper{}, nil
	case "checkbox":
		return &CheckboxMapper{}, nil
	case "date":
		return &DateMapper{}, nil
	}

	return nil, fmt.Errorf("unsupported property type: %v", propType)
}
