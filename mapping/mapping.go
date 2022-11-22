package mapping

import (
	"fmt"
	"reflect"

	"github.com/psyark/notionapi"
)

func Create(src interface{}, schema map[string]notionapi.Property) (*notionapi.CreatePageOptions, error) {
	return nil, nil
}

func Decode(page notionapi.Page, dst interface{}) error {
	dstType := reflect.TypeOf(dst)
	dstValue := reflect.ValueOf(dst)
	if dstType.Kind() != reflect.Pointer || dstType.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("dst must be *struct")
	}

	for i := 0; i < dstType.Elem().NumField(); i++ {
		f := dstType.Elem().Field(i)
		v := dstValue.Elem().Field(i)
		if tag, ok := f.Tag.Lookup("notion"); ok {
			pv := page.Properties.Get(tag)
			if pv == nil {
				return fmt.Errorf("unknown property id: %v", tag)
			}

			mapper, err := getMapper(pv.Type)
			if err != nil {
				return err
			}

			if err := mapper.RecordToObject(f, v, pv); err != nil {
				return err
			}
		}
	}

	return nil
}

func GetUpdatePageOptions(page notionapi.Page, src interface{}) (*notionapi.UpdatePageOptions, error) {
	srcType := reflect.TypeOf(src)
	srcValue := reflect.ValueOf(src)
	if srcType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("src must be struct")
	}

	opt := &notionapi.UpdatePageOptions{
		Properties: notionapi.PropertyValueMap{},
	}

	for i := 0; i < srcType.NumField(); i++ {
		f := srcType.Field(i)
		v := srcValue.Field(i)
		if tag, ok := f.Tag.Lookup("notion"); ok {
			pv := page.Properties.Get(tag)
			if pv == nil {
				return nil, fmt.Errorf("unknown property: %v", tag)
			}

			if pv2, err := compareField(f, v, pv); err != nil {
				return nil, err
			} else if pv2 != nil {
				opt.Properties[tag] = *pv2
			}
		}
	}

	return opt, nil
}

func compareField(field reflect.StructField, value reflect.Value, pv *notionapi.PropertyValue) (*notionapi.PropertyValue, error) {
	return nil, nil
}
