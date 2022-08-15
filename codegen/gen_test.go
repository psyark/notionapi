package codegen

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/dave/jennifer/jen"
)

func TestDatabase(t *testing.T) {
	url := "https://developers.notion.com/reference/database"
	builder := NewBuilder().Add(CommentWithBreak(url))

	err := Parse(url, func(title, desc string, props []DocProp) error {
		builder.AddClass("Database", desc).AddDocProps(props...)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := builder.Build("../types.database.go"); err != nil {
		t.Fatal(err)
	}
}

func TestPage(t *testing.T) {
	url := "https://developers.notion.com/reference/page"
	builder := NewBuilder().Add(CommentWithBreak(url))

	err := Parse(url, func(title, desc string, props []DocProp) error {
		builder.AddClass("Page", desc).AddDocProps(props...)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := builder.Build("../types.page.go"); err != nil {
		t.Fatal(err)
	}
}

func TestRichText(t *testing.T) {
	url := "https://developers.notion.com/reference/rich-text"
	builder := NewBuilder().Add(CommentWithBreak(url))

	err := Parse(url, func(title, desc string, props []DocProp) error {
		if strings.HasSuffix(title, " mentions") {
			builder.GetClass("Mention").AddField(Comment(desc)).AddDocProps(props...)
			return nil
		}

		switch title {
		case "All rich text":
			richText := builder.AddClass("RichText", desc)
			for _, dp := range props {
				p := dp.Property()
				if p.Name == "href" {
					p.Type = jen.Op("*").String()
					p.OmitEmpty = false
				}
				richText.AddField(p)
			}
		case "Annotations":
			builder.AddClass("Annotations", desc).AddDocProps(props...)
		case "Link objects":
			builder.AddClass("Link", desc).AddDocProps(props...)
		case "Text objects", "Mention objects", "Equation objects":
			className := strings.TrimSuffix(title, " objects")
			tagName := strings.ToLower(className)
			builder.AddClass(className, desc).AddDocProps(props...)
			builder.GetClass("RichText").AddConfiguration(tagName, className, desc)
		default:
			return fmt.Errorf("unknown title: %v", title)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := builder.Build("../types.richtext.go"); err != nil {
		t.Fatal(err)
	}
}

func TestProperty(t *testing.T) {
	configRegex := regexp.MustCompile(` (contain the following|have no additional) configuration within the (\w+) property`)

	url := "https://developers.notion.com/reference/property-object"
	builder := NewBuilder().Add(CommentWithBreak(url))

	err := Parse(url, func(title, desc string, props []DocProp) error {

		if strings.HasSuffix(title, " configuration") {
			desc = strings.ReplaceAll(desc, " ", " ")
			match := configRegex.FindStringSubmatch(desc)
			if len(match) == 0 {
				panic(desc)
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
		case "Multi-select options":
		default:
			return fmt.Errorf("unknown title: %v", title)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := builder.Build("../types.property.go"); err != nil {
		t.Fatal(err)
	}
}

func TestParent(t *testing.T) {
	url := "https://developers.notion.com/reference/parent-object"
	builder := NewBuilder().Add(CommentWithBreak(url))
	parent := builder.AddClass("Parent", `Pages, databases, and blocks are either located inside other pages, databases, and blocks, or are located at the top level of a workspace. This location is known as the "parent". Parent information is represented by a consistent parent object throughout the API.`)
	parent.AddField(&Property{Name: "type", Type: jen.String()})

	err := Parse(url, func(title, desc string, props []DocProp) error {
		fields := []Coder{}
		for _, dp := range props {
			if dp.Name != "type" {
				p := dp.Property()
				p.OmitEmpty = true
				fields = append(fields, p)
			}
		}
		parent.AddField(fields...)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := builder.Build("../types.parent.go"); err != nil {
		t.Fatal(err)
	}
}

func TestUser(t *testing.T) {
	url := "https://developers.notion.com/reference/user"
	builder := NewBuilder().Add(CommentWithBreak(url))

	err := Parse(url, func(title, desc string, props []DocProp) error {
		switch title {
		case "Where user objects appear in the API":
			builder.AddClass("User", desc)
		case "All users":
			builder.GetClass("User").AddField(Comment(desc)).AddDocProps(props...)
		case "People":
			person := builder.AddClass("Person", desc)
			for _, p := range props {
				if p.Name == "person" {
					builder.GetClass("User").AddField(Comment(p.Description)).AddConfiguration(p.Name, "Person", "")
				} else if strings.HasPrefix(p.Name, "person.") {
					p.Name = strings.TrimPrefix(p.Name, "person.")
					person.AddDocProps(p)
				} else {
					return fmt.Errorf("unknown property: %v", p.Name)
				}
			}
		case "Bots":
			if props[1].Name == "" && props[2].Name == "owner.type" && props[4].Name == "owner.type" {
				props = append(props[0:1], props[3:]...)
			} else {
				panic("ドキュメントが訂正されたようなので修正ロジックを見直してください")
			}

			for _, dp := range props {
				if dp.Name == "bot" {
					builder.GetClass("User").AddField(Comment(dp.Description)).AddConfiguration(dp.Name, "Bot", "")
					builder.AddClass("Bot", desc)
				} else if dp.Name == "owner" {
					builder.GetClass("Bot").AddField(&Property{
						Name:        dp.Name,
						Type:        jen.Op("*").Id("Owner"),
						Description: dp.Description,
					})
					builder.AddClass("Owner", desc)
				} else if strings.HasPrefix(dp.Name, "owner.") {
					dp.Name = strings.TrimPrefix(dp.Name, "owner.")
					p := dp.Property()
					builder.GetClass("Owner").AddField(p)
				} else {
					return fmt.Errorf("unknown property: %v", dp.Name)
				}
			}
		default:
			return fmt.Errorf("unknown title: %v", title)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := builder.Build("../types.user.go"); err != nil {
		t.Fatal(err)
	}
}

func TestPropertyValue(t *testing.T) {
	url := "https://developers.notion.com/reference/property-value-object"
	builder := NewBuilder().Add(CommentWithBreak(url))

	err := Parse(url, func(title, desc string, props []DocProp) error {
		if title == "All property values" {
			builder.AddClass("PropertyValue", desc).AddDocProps(props...)
		} else {
			builder.GetClass("PropertyValue").AddField(Comment(title + ": " + desc)) // TODO
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := builder.Build("../types.propertyvalue.go"); err != nil {
		t.Fatal(err)
	}
}

func TestFile(t *testing.T) {
	url := "https://developers.notion.com/reference/file-object"
	builder := NewBuilder().Add(CommentWithBreak(url))

	err := Parse(url, func(title, desc string, props []DocProp) error {
		switch title {
		case "All file objects":
			for i := range props {
				if props[i].Name == "{type}" {
					props[i].Name = "type"
				}
			}
			builder.AddClass("File", desc).AddDocProps(props...)
		case "Externally hosted files vs. Files hosted by Notion":
			builder.GetClass("File").Comment += fmt.Sprintf("\n\n%v\n%v", title, desc)
		default:
			match := regexp.MustCompile(`have a type of "(\w+)"`).FindStringSubmatch(desc)
			if len(match) != 0 {
				className := getName(strings.ReplaceAll(strings.TrimSuffix(title, " objects"), " ", "_"))
				builder.AddClass(className, title+"\n"+desc).AddDocProps(props...)
				builder.GetClass("File").AddConfiguration(match[1], className, desc)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := builder.Build("../types.file.go"); err != nil {
		t.Fatal(err)
	}
}

func TestEmoji(t *testing.T) {
	url := "https://developers.notion.com/reference/emoji-object"
	builder := NewBuilder().Add(CommentWithBreak(url))

	err := Parse(url, func(title, desc string, props []DocProp) error {
		builder.AddClass("Emoji", desc).AddDocProps(props...)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := builder.Build("../types.emoji.go"); err != nil {
		t.Fatal(err)
	}
}
