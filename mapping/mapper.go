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
		return &titleMapper{propMapper}, nil
	case "rich_text":
		return &richTextMapper{propMapper}, nil
	case "email":
		return &emailMapper{propMapper}, nil
	case "url":
		return &urlMapper{propMapper}, nil
	case "phone_number":
		return &phoneNumberMapper{propMapper}, nil
	case "number":
		return &numberMapper{propMapper}, nil
	case "checkbox":
		return &checkboxMapper{propMapper}, nil
	case "date":
		return &dateMapper{propMapper}, nil
	case "relation":
		return &relationMapper{propMapper}, nil
	}

	return nil, fmt.Errorf("unsupported property type: %v", propType)
}
