package codegen

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dave/jennifer/jen"
)

func BuildPropertyItem() error {
	url := "https://developers.notion.com/reference/property-item-object"
	builder := NewBuilder().Add(CommentWithBreak(url))

	descRegex := regexp.MustCompile(`contain (.+) within the (\w+) property`)

	{
		code := jen.Type().Id("PropertyItem").Interface(
			jen.Id("getCommon").Params().Op("*").Id("PropertyItemCommon"),
		).Line()
		builder.Add(RawCoder{Value: code})
	}

	cases := []jen.Code{}

	err := Parse(url, func(title, desc string, props []DocProp) error {
		if title == "All property items" {
			object := builder.AddClass("PropertyItemCommon", desc)
			for _, dp := range props {
				switch dp.Name {
				case "next_url": // Only present in paginated property values
					object.AddField(RawCoder{jen.Id("NextURL").Op("*").String().Tag(map[string]string{"json": "-"}).Comment(dp.Description)})
				default:
					object.AddDocProps(dp)
				}
			}
			code := jen.Func().Params(jen.Id("i").Op("*").Id("PropertyItemCommon")).Id("getCommon").Params().Op("*").Id("PropertyItemCommon").Block(
				jen.Return().Id("i"),
			)
			builder.Add(RawCoder{Value: code})
		} else if title == "Paginated property values" {
		} else if title == "Multi-select option values" {
			builder.AddClass(getName(title), desc).AddDocProps(props...)
		} else if strings.HasSuffix(title, " property values") {
			typeName := strings.TrimSuffix(title, " property values")
			typeName = strings.ReplaceAll(typeName, "-", "_")
			typeName = strings.ReplaceAll(typeName, " ", "_")
			typeName = strings.ToLower(typeName)

			name := getName(strings.Replace(title, " property values", " property item", 1))
			builder.AddClass(name, desc).AddField(AnonymousField("PropertyItemCommon"))

			{
				code := jen.Case(jen.Lit(typeName)).Return().List(
					jen.Op("&").Id(name).Block(),
					jen.Nil(),
				)
				cases = append(cases, code)
			}

			match := descRegex.FindStringSubmatch(desc)
			if len(match) != 0 {
				if match[1] == "the following data" {
					builder.GetClass(name).AddField(
						Property{
							Name: match[2],
							Type: jen.Id(name + "Data"),
						},
					)
					builder.AddClass(name+"Data", "").AddDocProps(props...)
				} else {
					prop := Property{Name: match[2], Description: desc}
					switch match[1] {
					case "a boolean":
						prop.Type = jen.Op("*").Bool()
					case "a string", "an optional string", "a non-empty string":
						prop.Type = jen.String()
					case "a number", "an optional number":
						prop.Type = jen.Float64()
					case "a user object":
						prop.Type = jen.Id("User")
					case "a date property value", "an optional date property value":
						prop.Type = jen.Id("DatePropertyItem")
					case "an array of rich text objects":
						prop.Type = jen.Id("RichText") // ignore "an array of". See https://developers.notion.com/reference/property-item-object#title-property-values
					case "an array of user objects":
						prop.Type = jen.Id("User") // ignore "an array of". See https://developers.notion.com/reference/property-item-object#title-property-values
					case "an array of file references":
						prop.Type = jen.Index().Id("File")
					case "an array of property_item objects":
						prop.Type = jen.Index().Id("DatePropertyItem")
					case "an array of multi-select option values":
						prop.Type = jen.Index().Id("MultiSelectOptionValues")
					case "an array of relation property items with a pagereferences":
						prop.Type = jen.Id("PageReference") // ignore "an array of". See https://developers.notion.com/reference/property-item-object#title-property-values
					default:
						fmt.Println(match[1])
					}
					if prop.Type != nil {
						builder.GetClass(name).AddField(prop)
					}
				}
			} else {
				// fmt.Println(desc)
			}
		} else {
			panic(title)
		}
		return nil
	})
	if err != nil {
		return err
	}

	{
		cases = append(cases, jen.Default().Return(
			jen.List(
				jen.Nil(),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit("unsupported type: %v"),
					jen.Id("typeName"),
				),
			),
		))
		code := jen.Func().Id("createPropertyItem").Params(jen.Id("typeName").String()).Params(
			jen.Id("PropertyItem"),
			jen.Error(),
		).Block(
			jen.Switch(jen.Id("typeName")).Block(cases...),
		)
		builder.Add(RawCoder{Value: code})
	}

	return builder.Build("../types.propertyitem.go")
}
