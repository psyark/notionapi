package codegen

import (
	"regexp"
	"strings"

	"github.com/dave/jennifer/jen"
)

func BuildPropertyItem() error {
	url := "https://developers.notion.com/reference/property-item-object"
	builder := NewBuilder().Add(CommentWithBreak(url))

	descRegex := regexp.MustCompile(`contain (.+) within the (\w+) property`)

	err := Parse(url, func(title, desc string, props []DocProp) error {
		if title == "All property items" {
			object := builder.AddClass("PropertyItem", desc)
			for _, dp := range props {
				switch dp.Name {
				case "next_url": // Only present in paginated property values
					object.AddField(RawCoder{jen.Id("NextURL").Op("*").String().Tag(map[string]string{"json": "-"}).Comment(dp.Description)})
				default:
					object.AddDocProps(dp)
				}
			}
		} else if title == "Paginated property values" {
		} else if title == "Multi-select option values" {
			// ignore
		} else if title == "Rollup property values" || title == "Formula property values" {
			name := getName(strings.Replace(title, " property values", " property item data", 1))
			builder.GetClass("PropertyItem").AddField(
				Property{
					Name:         strings.ToLower(strings.TrimSuffix(title, " property values")),
					Type:         jen.Id(name),
					Description:  desc,
					TypeSpecific: true,
				},
			)
			builder.AddClass(name, "").AddDocProps(props...)
		} else if strings.HasSuffix(title, " formula property values") {
			match := descRegex.FindStringSubmatch(desc)
			p := Property{Name: match[2], Description: desc, TypeSpecific: true}
			switch match[1] {
			case "an optional string":
				p.Type = jen.Op("*").String()
			case "an optional number":
				p.Type = jen.Op("*").Float64()
			case "a boolean":
				p.Type = jen.Bool()
			case "an optional date property value":
				p.Type = jen.Id("DatePropertyItemData")
			default:
				panic(match[1])
			}
			builder.GetClass("FormulaPropertyItemData").AddField(p)
		} else if title == "Incomplete rollup property values" {
		} else if strings.HasSuffix(title, " rollup property values") {
			match := descRegex.FindStringSubmatch(desc)
			p := Property{Name: match[2], Type: jen.Interface(), Description: desc, TypeSpecific: true}
			switch match[1] {
			case "a number":
				p.Type = jen.Float64()
			case "a date property value":
				p.Type = jen.Id("DatePropertyItemData")
			case "an array of property_item objects":
				p.Name = "array"
				p.Type = jen.Index().Struct()
			default:
				panic(match[1])
			}
			builder.GetClass("RollupPropertyItemData").AddField(p)
		} else if strings.HasSuffix(title, " property values") {
			typeName := strings.TrimSuffix(title, " property values")
			typeName = strings.ReplaceAll(typeName, "-", "_")
			typeName = strings.ReplaceAll(typeName, " ", "_")
			typeName = strings.ToLower(typeName)

			// tagName := strings.ToLower(strings.TrimSuffix(title, " property values"))
			name := getName(strings.Replace(title, " property values", " property item", 1))

			match := descRegex.FindStringSubmatch(desc)
			if len(match) != 0 {
				if match[1] == "the following data" {
					builder.GetClass("PropertyItem").AddField(
						Property{
							Name:         match[2],
							Type:         jen.Id(name + "Data"),
							Description:  desc,
							TypeSpecific: true,
						},
					)
					builder.AddClass(name+"Data", "").AddDocProps(props...)
				} else {
					prop := Property{Name: match[2], Description: desc, TypeSpecific: true}
					switch match[1] {
					case "a boolean":
						prop.Type = jen.Op("*").Bool()
					case "a string", "a non-empty string":
						prop.Type = jen.String()
					case "a number":
						prop.Type = jen.Float64()
					case "a user object":
						prop.Type = jen.Id("UserPropertyValueData")
					case "an array of rich text objects":
						prop.Type = jen.Id("RichText") // ignore "an array of". See https://developers.notion.com/reference/property-item-object#title-property-values
					case "an array of user objects":
						prop.Type = jen.Id("UserPropertyValueData") // ignore "an array of". See https://developers.notion.com/reference/property-item-object#title-property-values
					case "an array of file references":
						prop.Type = jen.Index().Id("File")
					case "an array of multi-select option values":
						prop.Type = jen.Index().Id("SelectPropertyItemData")
					case "an array of relation property items with page references":
						prop.Type = jen.Id("PageReference") // ignore "an array of". See https://developers.notion.com/reference/property-item-object#title-property-values
					default:
						panic(match[1])
					}
					builder.GetClass("PropertyItem").AddField(prop)
				}
			} else {
				panic(title)
			}
		} else {
			panic(title)
		}
		return nil
	})
	if err != nil {
		return err
	}

	builder.GetClass("PropertyItem").AddField(Property{
		Name:         "status",
		Type:         jen.Id("SelectPropertyItemData"),
		Description:  "undocumented",
		TypeSpecific: true,
	})

	return builder.Build("../types.propertyitem.go")
}
