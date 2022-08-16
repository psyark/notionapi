package codegen

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dave/jennifer/jen"
)

var _ Coder = &DocProp{}

// DocProp はAPI Documentに記載されたプロパティのCoderです
type DocProp struct {
	Name        string
	Type        string
	Description string
}

func (dp DocProp) Property() *Property {
	typeCode, omitEmpty := dp.getType()
	return &Property{
		Name:        strings.TrimSuffix(strings.TrimSuffix(dp.Name, "*"), " (optional)"),
		Type:        typeCode,
		Description: dp.Description,
		OmitEmpty:   omitEmpty,
	}
}

func (dp DocProp) Code() jen.Code {
	return dp.Property().Code()
}

func (dp DocProp) getType() (jen.Code, bool) {
	switch dp.Type {
	case "null":
		return jen.Interface(), false
	case "boolean", "boolean (optional)":
		return jen.Bool(), false
	case "integer":
		return jen.Int64(), false
	case "string", "string enum", "string (enum)", "string (optional)", "string (optional enum)", "string (optional, enum)":
		return jen.String(), strings.Contains(dp.Name, "optional") || strings.Contains(dp.Type, "optional")
	case "string or null":
		return jen.Op("*").String(), false
	case "string (UUID)", "string (UUIDv4)":
		return jen.Id("UUIDString"), false
	case "string (ISO 8601 date time)", "string (ISO 8601 date and time)", "string (optional, ISO 8601 date and time)":
		return jen.Id("ISO8601String"), false
	case "Partial User":
		return jen.Op("*").Id("User"), false
	case "File Object or Emoji object", `File Object (only type of "external" is supported currently) or Emoji object`:
		return jen.Op("*").Id("FileOrEmoji"), false
	case "File Object", `File object (only type of "external" is supported currently)`:
		return jen.Op("*").Id("File"), false
	case "Synced From Object":
		return jen.Op("*").Id("SyncedFrom"), false
	case "array of block objects":
		return jen.Index().Id("Block"), false
	case "array of rich text objects", "array of Rich text object text objects":
		return jen.Index().Id("RichText"), false
	case "array of array of Rich text objects":
		return jen.Index().Index().Id("RichText"), false
	case "array of select option objects.", "array of multi-select option objects.":
		return jen.Index().Id("SelectOption"), false
	case "array of table_row block objects":
		return jen.Index().Interface(), false
	case "array of column_list block objects":
		return jen.Index().Interface(), false
	case "object", "object (optional)":
		switch dp.Name {
		case "annotations", "link", "parent", "user":
			return jen.Op("*").Id(getName(dp.Name)), false
		case "properties", "properties*":
			if strings.Contains(dp.Description, "Property object") {
				return jen.Map(jen.String()).Id("Property"), false
			} else if strings.Contains(dp.Description, "Property value object") {
				return jen.Map(jen.String()).Id("PropertyValue"), false
			}
			panic(dp.Description)
		default:
			return jen.Interface(), false
		}
	case "list":
		if dp.Description == "List of property_item objects." {
			return jen.Index().Id("PropertyItem"), false
		} else {
			panic(fmt.Errorf("getType: %v", dp.Type))
		}
	default:
		if strings.HasPrefix(dp.Type, "string (enum)") {
			return jen.String(), false
		} else if regexp.MustCompile(`^"\w+"$`).MatchString(dp.Type) {
			return jen.String(), false
		} else {
			panic(fmt.Errorf("getType: %v", dp.Type))
		}
	}
}
