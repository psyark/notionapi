package codegen

import (
	"regexp"
	"strings"
	"testing"

	"github.com/dave/jennifer/jen"
)

func TestPropertyValueObject(t *testing.T) {
	t.Parallel()

	url := "https://developers.notion.com/reference/property-value-object"
	builder := NewBuilder().Add(CommentWithBreak(url))

	descRegex := regexp.MustCompile(`contain (?:a |an )?(.+) within the (\w+) property`)

	sections, err := ParseObjectDoc(url)
	if err != nil {
		t.Fatal(err)
	}

	for _, section := range sections {
		heading := section.Heading
		desc := section.AllParagraphText()

		if heading == nil {
			builder.AddClass("PropertyValue", desc)
			continue
		}

		title := heading.Text
		switch title {
		case "All property values":
			obj := builder.GetClass("PropertyValue").AddField(Comment(desc))
			for _, param := range section.Parameters() {
				if param.Name == "type (optional)" {
					param.Name = "type"
				}

				opt := &PropertyOption{OmitEmpty: true} // RollupPropertyValueData.Array ではIDが無いため
				if err := obj.AddParams(opt, param); err != nil {
					t.Fatal(err)
				}
			}
		case "Select property values":
			builder.GetClass("PropertyValue").AddConfiguration("select", "SelectOption", desc)
		case "Status property values":
			builder.GetClass("PropertyValue").AddConfiguration("status", "StatusOption", desc)
		case "Multi-select property values":
			prop := &Property{Name: "multi_select", Type: jen.Index().Id("SelectOption"), Description: desc, TypeSpecific: true}
			builder.GetClass("PropertyValue").AddField(prop)
		case "Multi-select option values": // ignore
		case "Date property values":
			builder.GetClass("PropertyValue").AddConfiguration("date", "DateValue", desc)
			builder.AddClass("DateValue", desc).AddParams(nil, section.Parameters()...)
		default:
			if strings.HasSuffix(title, " formula property values") {
				match := descRegex.FindStringSubmatch(desc)
				param := ObjectDocParameter{Name: match[2], Type: match[1], Description: desc}
				opt := &PropertyOption{TypeSpecific: true}
				if err := builder.GetClass("FormulaPropertyValueData").AddParams(opt, param); err != nil {
					t.Fatal(err)
				}
			} else if strings.HasSuffix(title, " rollup property values") {
				match := descRegex.FindStringSubmatch(desc)
				switch match[1] {
				case "array of number, date, or string objects":
					// TODO: 確認
					prop := &Property{Name: "array", Type: jen.Index().Id("PropertyValue"), Description: desc, TypeSpecific: true}
					builder.GetClass("RollupPropertyValueData").AddField(prop)
				default:
					param := ObjectDocParameter{Name: match[2], Type: match[1], Description: desc}
					opt := &PropertyOption{TypeSpecific: true}
					if err := builder.GetClass("RollupPropertyValueData").AddParams(opt, param); err != nil {
						t.Fatal(err)
					}
				}
			} else if strings.HasSuffix(title, " property values") {
				if match := descRegex.FindStringSubmatch(desc); len(match) != 0 {
					param := ObjectDocParameter{Name: match[2], Type: match[1], Description: desc}

					switch match[1] {
					case "the following data":
						dataName := nfCamelCase.String(strings.TrimSuffix(title, "s")) + "Data"
						if err := builder.AddClass(dataName, desc).AddParams(nil, section.Parameters()...); err != nil {
							t.Fatal(err)
						}

						prop := &Property{Name: match[2], Type: jen.Op("*").Id(dataName), Description: desc, TypeSpecific: true}
						builder.GetClass("PropertyValue").AddField(prop)
					default:
						opt := &PropertyOption{TypeSpecific: true, Nullable: true}
						if err := builder.GetClass("PropertyValue").AddParams(opt, param); err != nil {
							t.Fatal(err)
						}
					}
				} else {
					prefix := strings.TrimSuffix(title, " property values")
					dataName := nfCamelCase.String(prefix) + "PropertyValueData"
					prop := &Property{Name: nf_snake_case.String(prefix), Type: jen.Id(dataName), Description: desc, TypeSpecific: true}
					if err := builder.AddClass(dataName, desc).AddParams(nil, section.Parameters()...); err != nil {
						t.Fatal(err)
					}
					if title == "Rollup property values" {
						builder.GetClass(dataName).AddField(&Property{Name: "type", Type: jen.String(), Description: "These objects contain a type key and a key corresponding with the value of type."})
						builder.GetClass(dataName).AddField(&Property{Name: "function", Type: jen.String(), Description: "undocumented"})
					}
					builder.GetClass("PropertyValue").AddField(prop)
				}
			} else {
				t.Error(title)
			}
		}
	}

	if err := builder.Build("../types.propertyvalue.go"); err != nil {
		t.Fatal(err)
	}
}
