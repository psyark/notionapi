package codegen

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
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

	for _, section := range sections {
		heading := section.Heading
		desc := section.ParagraphText()

		if heading == nil {
			continue
		}

		title := heading.Text

		if strings.HasSuffix(title, " configuration") {
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
				clsName := getName(strings.ReplaceAll(strings.ReplaceAll(title, "-", "_"), " ", "_"))

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
					builder.GetClass("Property").AddConfiguration(match[2], clsName, desc)
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
					builder.GetClass("Property").AddConfiguration(match[2], clsName, desc)
				}
			}
		} else {
			switch title {
			case "Database properties":
				if err := builder.AddClass("Property", desc).AddParams(nil, section.Parameters()...); err != nil {
					t.Error(err)
				}
			case "Select options":
				if err := builder.AddClass("SelectOption", desc).AddParams(nil, section.Parameters()...); err != nil {
					t.Error(err)
				}
			case "Status options":
				if err := builder.AddClass("StatusOption", desc).AddParams(nil, section.Parameters()...); err != nil {
					t.Error(err)
				}
			case "Status groups":
				if err := builder.AddClass("StatusGroup", desc).AddParams(nil, section.Parameters()...); err != nil {
					t.Error(err)
				}
			case "Multi-select options":
			default:
				t.Fatal(title)
			}
		}
	}

	if err := builder.Build("../types.property.go"); err != nil {
		t.Fatal(err)
	}
}
