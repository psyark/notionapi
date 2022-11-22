package codegen

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/dave/jennifer/jen"
)

func TestPropertyObject(t *testing.T) {
	t.Parallel()

	url := "https://developers.notion.com/reference/property-object"
	builder := NewBuilder().Add(CommentWithBreak(url))

	configRegex := regexp.MustCompile(` (contain the following|have no additional) configuration within the (\w+) property`)

	sections, err := ParseObjectDoc(url)
	if err != nil {
		t.Fatal(err)
	}

	property := builder.AddClass("Property", sections[0].AllParagraphText())

	builder.Add(RawCoder{jen.Type().Id("PropertyMap").Map(jen.String()).Id("Property")})
	builder.Add(RawCoder{jen.Func().Params(jen.Id("m").Id("PropertyMap")).Id("Get").Params(jen.Id("id").String()).Op("*").Id("Property").Block(
		jen.For().List(jen.Id("_"), jen.Id("pv")).Op(":=").Range().Id("m").Block(
			jen.If(jen.Id("pv").Dot("ID").Op("==").Id("id")).Block(
				jen.Return().Op("&").Id("pv"),
			),
		),
		jen.Return().Nil(),
	)})

	for _, section := range sections[1:] {
		title := section.Heading.Text
		desc := section.AllParagraphText()

		switch title {
		case "Database properties":
			if err := property.AddField(Comment(desc)).AddParams(nil, section.Parameters()...); err != nil {
				t.Error(err)
			}
			property.AddLine()

		case "Select options":
			opt := &PropertyOption{OmitEmpty: true} // UpdatePageのPropertyValueでは省略可能
			if err := builder.AddClass("SelectOption", desc).AddParams(opt, section.Parameters()...); err != nil {
				t.Error(err)
			}
		case "Status options":
			opt := &PropertyOption{OmitEmpty: true} // UpdatePageのPropertyValueでは省略可能
			if err := builder.AddClass("StatusOption", desc).AddParams(opt, section.Parameters()...); err != nil {
				t.Error(err)
			}
		case "Status groups":
			if err := builder.AddClass("StatusGroup", desc).AddParams(nil, section.Parameters()...); err != nil {
				t.Error(err)
			}
		case "Multi-select options": // ignore
		case "Title configuration",
			"Text configuration",
			"Number configuration",
			"Select configuration",
			"Status configuration",
			"Multi-select configuration",
			"Date configuration",
			"People configuration",
			"Files configuration",
			"Checkbox configuration",
			"URL configuration",
			"Email configuration",
			"Phone number configuration",
			"Formula configuration",
			"Relation configuration",
			"Single property relation configuration",
			"Dual property relation configuration",
			"Rollup configuration",
			"Created time configuration",
			"Created by configuration",
			"Last edited time configuration",
			"Last edited by configuration":

			desc = strings.ReplaceAll(desc, " ", " ")

			if title == "Status configuration" {
				if desc == "" {
					desc = "Status database property objects contain the following configuration within the status property:"
				} else {
					panic(fmt.Sprintf("title=%v, desc=%v", title, desc))
				}
			}

			match := configRegex.FindStringSubmatch(desc)
			if len(match) == 0 {
				panic(fmt.Sprintf("title=%v, desc=%v", title, desc))
			}

			if match[1] == "have no additional" {
				if title == "Single property relation configuration" {
					containerName := "RelationConfiguration"
					builder.GetClass(containerName).AddConfiguration(match[2], "", desc)
				} else {
					containerName := "Property"
					builder.GetClass(containerName).AddConfiguration(match[2], "", desc)
				}
			} else {
				clsName := nfCamelCase.String(title)

				switch title {
				case "Relation configuration":
					object := builder.AddClass(clsName, desc)
					for _, dp := range section.Parameters() {
						if dp.Name == "type" {
							dp.Type = "string" // string (optional enum) -> string
						}
						if err := object.AddParams(nil, dp); err != nil {
							t.Fatal(err)
						}
					}
					property.AddConfiguration(match[2], clsName, desc)
				case "Dual property relation configuration":
					obj := builder.AddClass(clsName, desc)
					for _, param := range section.Parameters() {
						if param.Name == "single_property" {
							param.Name = "synced_property_id"
						}
						if err := obj.AddParams(nil, param); err != nil {
							t.Fatal(err)
						}
					}
					builder.GetClass("RelationConfiguration").AddConfiguration(match[2], clsName, desc)
				default:
					if err := builder.AddClass(clsName, desc).AddParams(nil, section.Parameters()...); err != nil {
						t.Fatal(err)
					}
					property.AddConfiguration(match[2], clsName, desc)
				}
			}

		default:
			t.Fatal(title)
		}
	}

	if err := builder.Build("../types.property.go"); err != nil {
		t.Fatal(err)
	}
}
