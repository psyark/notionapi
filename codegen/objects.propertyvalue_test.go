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

	property_value := builder.AddClass("PropertyValue", sections[0].AllParagraphText())

	for _, section := range sections[1:] {
		title := section.Heading.Text
		desc := section.AllParagraphText()

		switch title {
		case "All property values":
			obj := property_value.AddField(Comment(desc))
			for _, param := range section.Parameters() {
				if param.Name == "type (optional)" {
					param.Name = "type"
				}

				opt := &PropertyOption{OmitEmpty: true} // RollupPropertyValueData.Array ではIDが無いため
				if err := obj.AddParams(opt, param); err != nil {
					t.Fatal(err)
				}
			}

			// TODO: 2022-09-30 にtype="relation"に突然出現したため仮対応。要らなくなったら外す
			property_value.AddField(&Property{Name: "has_more", Type: jen.Op("*").Bool(), OmitEmpty: true, Description: "undocumented"})
		case "Select property values", "Status property values":
			prefix := strings.TrimSuffix(title, " property values")
			property_value.AddConfiguration(nf_snake_case.String(prefix), prefix+"Option", desc)
		case "Multi-select property values":
			prop := &Property{Name: "multi_select", Type: jen.Index().Id("SelectOption"), Description: desc, TypeSpecific: true}
			property_value.AddField(prop)
		case "Multi-select option values": // ignore
		case "Date property values":
			property_value.AddConfiguration("date", "DateValue", desc)
			builder.AddClass("DateValue", desc).AddParams(nil, section.Parameters()...)
		case "String formula property values", "Number formula property values", "Boolean formula property values", "Date formula property values":
			match := descRegex.FindStringSubmatch(desc)
			param := ObjectDocParameter{Name: match[2], Type: match[1], Description: desc}
			opt := &PropertyOption{TypeSpecific: true}
			if err := builder.GetClass("FormulaPropertyValueData").AddParams(opt, param); err != nil {
				t.Fatal(err)
			}
		case "String rollup property values", "Number rollup property values", "Date rollup property values", "Array rollup property values":
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
		case
			"Title property values", "Rich Text property values", "Number property values", "Formula property values",
			"Relation property values", "Rollup property values", "People property values", "Files property values",
			"Checkbox property values", "URL property values", "Email property values", "Phone number property values",
			"Created time property values", "Created by property values", "Last edited time property values", "Last edited by property values":
			if match := descRegex.FindStringSubmatch(desc); len(match) != 0 {
				param := ObjectDocParameter{Name: match[2], Type: match[1], Description: desc}

				switch match[1] {
				case "the following data":
					dataName := nfCamelCase.String(strings.TrimSuffix(title, "s")) + "Data"
					if err := builder.AddClass(dataName, desc).AddParams(nil, section.Parameters()...); err != nil {
						t.Fatal(err)
					}

					prop := &Property{Name: match[2], Type: jen.Op("*").Id(dataName), Description: desc, TypeSpecific: true}
					property_value.AddField(prop)
				default:
					opt := &PropertyOption{TypeSpecific: true, Nullable: true}
					if err := property_value.AddParams(opt, param); err != nil {
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
				property_value.AddField(prop)
			}
		default:
			t.Error(title)
		}
	}

	if err := builder.Build("../types.propertyvalue.go"); err != nil {
		t.Fatal(err)
	}
}
