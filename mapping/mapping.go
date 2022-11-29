package mapping

import (
	"fmt"
	"reflect"

	"github.com/psyark/notionapi"
)

func DecodePage(dst interface{}, page notionapi.Page) error {
	dstType := reflect.TypeOf(dst)
	dstValue := reflect.ValueOf(dst)
	if dstType.Kind() != reflect.Pointer || dstType.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("dst must be *struct")
	}

	for i := 0; i < dstType.Elem().NumField(); i++ {
		f := dstType.Elem().Field(i)
		v := dstValue.Elem().Field(i)
		if tag := parseTag(f); tag != nil {
			mapper, err := getMapperByPage(tag, page)
			if err != nil {
				return err
			}

			if err := mapper.decodePage(v, page); err != nil {
				return err
			}
		}
	}

	return nil
}

func CreatePageFrom(src interface{}, database *notionapi.Database) (*notionapi.CreatePageOptions, error) {
	srcType := reflect.TypeOf(src)
	srcValue := reflect.ValueOf(src)
	if srcType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("src must be struct")
	}

	opt := &notionapi.CreatePageOptions{
		Parent:     &notionapi.Parent{Type: "database_id", DatabaseID: database.ID},
		Properties: notionapi.PropertyValueMap{},
	}

	for i := 0; i < srcType.NumField(); i++ {
		f := srcType.Field(i)
		v := srcValue.Field(i)
		if tag := parseTag(f); tag != nil {
			mapper, err := getMapperByDB(tag, *database)
			if err != nil {
				return nil, err
			}

			if err := mapper.createPageFrom(v, opt); err != nil {
				return nil, err
			}
		}
	}

	if len(opt.Properties) != 0 || opt.Icon != nil || opt.Cover != nil {
		return opt, nil
	} else {
		return nil, nil
	}
}

func UpdatePageFrom(src interface{}, page notionapi.Page) (*notionapi.UpdatePageOptions, error) {
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
		if tag := parseTag(f); tag != nil {
			mapper, err := getMapperByPage(tag, page)
			if err != nil {
				return nil, err
			}

			if err := mapper.updatePageFrom(v, page, opt); err != nil {
				return nil, err
			}
		}
	}

	if len(opt.Properties) != 0 || opt.Icon != nil || opt.Cover != nil {
		return opt, nil
	} else {
		return nil, nil
	}
}
