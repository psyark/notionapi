package codegen

import (
	"regexp"
	"strings"

	"github.com/dave/jennifer/jen"
)

func BuildPropertyValue() error {
	url := "https://developers.notion.com/reference/property-value-object"
	builder := NewBuilder().Add(CommentWithBreak(url))

	descRegex := regexp.MustCompile(`contain (.+) within the (\w+) property`)

	err := Parse(url, func(title, desc string, props []DocProp) error {
		if title == "All property values" {
			obj := builder.AddClass("PropertyValue", desc)
			for _, dp := range props {
				if dp.Name == "id" {
					p := dp.Property()
					p.OmitEmpty = true // RollupPropertyValueData.Array ではIDが無いため
					obj.AddField(p)
				} else {
					obj.AddDocProps(dp)
				}
			}
		} else if title == "Multi-select option values" {
			// ignore
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
			builder.GetClass("FormulaPropertyValueData").AddField(p)
		} else if strings.HasSuffix(title, " rollup property values") {
			match := descRegex.FindStringSubmatch(desc)
			p := Property{Name: match[2], Type: jen.Interface(), Description: desc, TypeSpecific: true}
			switch match[1] {
			case "an optional string":
				p.Type = jen.String()
			case "a number":
				p.Type = jen.Float64()
			case "a date property value":
				p.Type = jen.Id("DatePropertyItemData")
			case "an array of number, date, or string objects":
				p.Name = "array"
				p.Type = jen.Index().Id("PropertyValue") // 確認する
			default:
				panic(match[1])
			}
			builder.GetClass("RollupPropertyValueData").AddField(p)
		} else if strings.HasSuffix(title, " property values") {
			match := descRegex.FindStringSubmatch(desc)
			prop := Property{Description: desc, TypeSpecific: true}

			if len(match) != 0 {
				prop.Name = match[2]
				switch match[1] {
				case "an array of rich text objects":
					prop.Type = jen.Index().Id("RichText")
				case "a number":
					prop.Type = jen.Float64()
				case "the following data":
					dataName := getName(strings.TrimSuffix(title, " property values") + " property value data")
					prop.Type = jen.Id(dataName)
					dataObj := builder.AddClass(dataName, desc)
					for _, dp := range props {
						p := dp.Property()
						p.OmitEmpty = false // 少なくとも DatePropertyValueData のため
						dataObj.AddField(p)
					}
				case "an array of multi-select option values":
					prop.Type = jen.Index().Id("SelectPropertyValueData")
				case "an array of user objects":
					prop.Type = jen.Index().Id("User")
				case "an array of file references":
					prop.Type = jen.Index().Id("File")
				case "a boolean":
					prop.Type = jen.Bool()
				case "a string", "a non-empty string":
					prop.Type = jen.String()
				case "a user object":
					prop.Type = jen.Id("User")
				case "an array of page references":
					prop.Type = jen.Index().Id("PageReference")
				default:
					panic(match[1])
				}
			} else {
				name := strings.ToLower(strings.TrimSuffix(title, " property values"))
				dataName := getName(name) + "PropertyValueData"
				prop.Name = strings.ToLower(name)
				prop.Type = jen.Id(dataName)
				builder.AddClass(dataName, desc).AddDocProps(props...)
				if title == "Rollup property values" {
					builder.GetClass(dataName).AddField(Property{Name: "type", Type: jen.String(), Description: "These objects contain a type key and a key corresponding with the value of type."})
					builder.GetClass(dataName).AddField(Property{Name: "function", Type: jen.String(), Description: "undocumented"})
				}
			}

			builder.GetClass("PropertyValue").AddField(prop)
		} else {
			panic(title)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return builder.Build("../types.propertyvalue.go")
}
