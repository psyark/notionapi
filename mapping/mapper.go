package mapping

import (
	"fmt"
	"reflect"

	"github.com/psyark/notionapi"
)

type mapper interface {
	decodePage(value reflect.Value, page notionapi.Page) error
	createPageFrom(value reflect.Value, dst *notionapi.CreatePageOptions) error
	updatePageFrom(value reflect.Value, page notionapi.Page, dst *notionapi.UpdatePageOptions) error
}

func getMapperByPage(tag *tagInfo, page notionapi.Page) (mapper, error) {
	if tag.icon {
		return &iconMapper{}, nil
	} else if tag.cover {
		return &coverMapper{}, nil
	}

	pv := page.Properties.Get(tag.name)
	if pv == nil {
		return nil, fmt.Errorf("unknown property id: %v", tag)
	}

	return getPropertyMapper(tag, tag.name, pv.Type)
}

func getMapperByDB(tag *tagInfo, db notionapi.Database) (mapper, error) {
	if tag.icon {
		return &iconMapper{}, nil
	} else if tag.cover {
		return &coverMapper{}, nil
	}

	prop := db.Properties.Get(tag.name)
	if prop == nil {
		return nil, fmt.Errorf("unknown property id: %v", tag)
	}

	return getPropertyMapper(tag, tag.name, prop.Type)
}

func getPropertyMapper(tag *tagInfo, propID string, propType string) (mapper, error) {
	propMapper := propertyMapper{propID}

	switch propType {
	case "title":
		return &richTextArrayMapper{
			getRTA: func(page notionapi.Page) notionapi.RichTextArray { return page.Properties.Get(propID).Title },
			setRTA: func(rta notionapi.RichTextArray, dst notionapi.PropertyValueMap) {
				dst[propID] = notionapi.PropertyValue{Type: "title", Title: rta}
			},
		}, nil
	case "rich_text":
		return &richTextArrayMapper{
			getRTA: func(page notionapi.Page) notionapi.RichTextArray { return page.Properties.Get(propID).RichText },
			setRTA: func(rta notionapi.RichTextArray, dst notionapi.PropertyValueMap) {
				dst[propID] = notionapi.PropertyValue{Type: "rich_text", RichText: rta}
			},
		}, nil
	case "email":
		return &stringPtrMapper{
			getStrPtr: func(page notionapi.Page) *string { return page.Properties.Get(propID).Email },
			setStrPtr: func(strPtr *string, dst notionapi.PropertyValueMap) {
				dst[propID] = notionapi.PropertyValue{Type: "email", Email: strPtr}
			},
		}, nil
	case "url":
		return &stringPtrMapper{
			getStrPtr: func(page notionapi.Page) *string { return page.Properties.Get(propID).URL },
			setStrPtr: func(strPtr *string, dst notionapi.PropertyValueMap) {
				dst[propID] = notionapi.PropertyValue{Type: "url", URL: strPtr}
			},
		}, nil
	case "phone_number":
		return &stringPtrMapper{
			getStrPtr: func(page notionapi.Page) *string { return page.Properties.Get(propID).PhoneNumber },
			setStrPtr: func(strPtr *string, dst notionapi.PropertyValueMap) {
				dst[propID] = notionapi.PropertyValue{Type: "phone_number", PhoneNumber: strPtr}
			},
		}, nil
	case "number":
		return &numberMapper{propMapper}, nil
	case "checkbox":
		return &checkboxMapper{propMapper}, nil
	case "date":
		return &dateMapper{propMapper}, nil
	case "relation":
		return &relationMapper{propMapper}, nil
	case "select":
		return &selectMapper{propMapper}, nil
	case "multi_select":
		return &multiSelectMapper{propMapper}, nil
	}

	return nil, fmt.Errorf("unsupported property type: %v", propType)
}
