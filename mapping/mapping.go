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

			// if err := decodeField(f, v, pv); err != nil {
			// 	return err
			// }
		}
	}

	return nil
}

func decodeField(field reflect.StructField, value reflect.Value, pv *notionapi.PropertyValue) error {
	type combi struct {
		propType  string
		fieldKind reflect.Kind
	}

	switch (combi{pv.Type, field.Type.Kind()}) {
	case combi{"title", reflect.Interface}, combi{"title", reflect.String}:
		value.Set(reflect.ValueOf(pv.Title.PlainText()))
	case combi{"rich_text", reflect.Interface}, combi{"rich_text", reflect.String}:
		value.Set(reflect.ValueOf(pv.RichText.PlainText()))
	case combi{"checkbox", reflect.Interface}, combi{"checkbox", reflect.Bool}:
		value.Set(reflect.ValueOf(pv.Checkbox))
	case combi{"email", reflect.Interface}, combi{"email", reflect.String}:
		if pv.Email == nil {
			value.Set(reflect.ValueOf(""))
		} else {
			value.Set(reflect.ValueOf(*pv.Email))
		}
	case combi{"url", reflect.Interface}, combi{"url", reflect.String}:
		if pv.URL == nil {
			value.Set(reflect.ValueOf(""))
		} else {
			value.Set(reflect.ValueOf(*pv.URL))
		}
	case combi{"phone_number", reflect.Interface}, combi{"phone_number", reflect.String}:
		if pv.PhoneNumber == nil {
			value.Set(reflect.ValueOf(""))
		} else {
			value.Set(reflect.ValueOf(*pv.PhoneNumber))
		}
	case combi{"number", reflect.Interface}, combi{"number", reflect.Float64}:
		if pv.Number == nil {
			value.Set(reflect.ValueOf(0.0))
		} else {
			value.Set(reflect.ValueOf(*pv.Number))
		}
	case combi{"number", reflect.Int}:
		if pv.Number == nil {
			value.SetInt(0)
		} else {
			value.SetInt(int64(*pv.Number))
		}
	case combi{"date", reflect.Interface}, combi{"date", reflect.String}:
		if pv.Date == nil {
			value.Set(reflect.ValueOf(""))
		} else {
			value.Set(reflect.ValueOf(pv.Date.Start))
		}
	default:
		return fmt.Errorf("unsupported combination: %v + %v", pv.Type, field.Type)
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
