package codegen

import (
	"fmt"
	"regexp"
	"strings"
)

func BuildProperty() error {
	configRegex := regexp.MustCompile(` (contain the following|have no additional) configuration within the (\w+) property`)

	url := "https://developers.notion.com/reference/property-object"
	builder := NewBuilder().Add(CommentWithBreak(url))

	err := Parse(url, func(title, desc string, props []DocProp) error {

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
				containerName := "Property"
				if title == "Single property relation configuration" {
					containerName = "RelationConfiguration"
				}
				builder.GetClass(containerName).AddConfiguration(match[2], "", desc)
			} else {
				clsName := getName(strings.ReplaceAll(strings.ReplaceAll(title, "-", "_"), " ", "_"))

				switch title {
				case "Relation configuration":
					object := builder.AddClass(clsName, desc)
					for _, dp := range props {
						if dp.Name != "synced_property_id" && dp.Name != "synced_property_name" {
							object.AddField(dp)
						}
					}
					builder.GetClass("Property").AddConfiguration(match[2], clsName, desc)
				case "Dual property relation configuration":
					for i := range props {
						if props[i].Name == "single_property" {
							props[i].Name = "synced_property_id"
						}
					}
					builder.AddClass(clsName, desc).AddDocProps(props...)
					builder.GetClass("RelationConfiguration").AddConfiguration(match[2], clsName, desc)
				default:
					builder.AddClass(clsName, desc).AddDocProps(props...)
					builder.GetClass("Property").AddConfiguration(match[2], clsName, desc)
				}

			}
			return nil
		}

		switch title {
		case "Database properties":
			builder.AddClass("Property", desc).AddDocProps(props...)
		case "Select options":
			builder.AddClass("SelectOption", desc).AddDocProps(props...)
		case "Status options":
			builder.AddClass("StatusOption", desc).AddDocProps(props...)
		case "Status groups":
			builder.AddClass("StatusGroup", desc).AddDocProps(props...)
		case "Multi-select options":
		default:
			return fmt.Errorf("unknown title: %v", title)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return builder.Build("../types.property.go")
}
