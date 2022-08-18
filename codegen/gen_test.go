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

func TestPagination(t *testing.T) {
	url := "https://developers.notion.com/reference/pagination"
	builder := NewBuilder().Add(CommentWithBreak(url))

	err := Parse(url, func(title, desc string, props []DocProp) error {
		if title == "Responses" {
			object := builder.AddClass("Pagination", desc)
			for _, p := range props {
				switch p.Name {
				case "results", "{type}": // ignore
				case "next_cursor":
					prop := p.Property()
					prop.Type = jen.Op("*").String() // next_cursor might be null.
					object.AddField(prop)
				case "type":
					match := regexp.MustCompile(` "(\w+)"`).FindAllStringSubmatch(p.Description, -1)
					for _, m := range match {
						name := getName(m[1])
						if name == "PropertyItem" {
							builder.AddClass(name+"Pagination", "").AddField(
								AnonymousField("Pagination"),
								Property{Name: "results", Type: jen.Index().Id("PropertyItemMarshaler")},
								Property{Name: m[1], Type: jen.Id("PropertyItemMarshaler")},
							)
						} else {
							builder.AddClass(name+"Pagination", "").AddField(
								AnonymousField("Pagination"),
								Property{Name: "results", Type: jen.Index().Id(name)},
								Property{Name: m[1], Type: jen.Struct()},
							)
						}
					}
					fallthrough
				default:
					object.AddDocProps(p)
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := builder.Build("../types.pagination.go"); err != nil {
		t.Fatal(err)
	}
}

func TestBlock(t *testing.T) {
	url := "https://developers.notion.com/reference/block"
	builder := NewBuilder().Add(CommentWithBreak(url))

	err := Parse(url, func(title, desc string, props []DocProp) error {
		switch title {
		case "Block object keys":
			object := builder.AddClass("Block", desc)
			for _, p := range props {
				if p.Name != "{type}" {
					object.AddDocProps(p)
				}
			}
		case "Block Type Object":
		case "Column List and Column Blocks":
			builder.AddClass("ColumnListBlocks", desc)
			builder.AddClass("ColumnBlocks", desc)
			builder.GetClass("Block").AddConfiguration("column_list", "ColumnListBlocks", desc)
			builder.GetClass("Block").AddConfiguration("column", "ColumnBlocks", desc)
		case "Synced Block blocks":
			builder.AddClass("SyncedBlockBlocks", desc).AddDocProps(props[2])
			builder.GetClass("Block").AddConfiguration("synced_block", "SyncedBlockBlocks", desc)
			builder.AddClass("SyncedFrom", desc).AddDocProps(props[3:]...)
		default:
			if strings.HasSuffix(title, " Blocks") || strings.HasSuffix(title, " blocks") {
				tagName := strings.ReplaceAll(strings.TrimSuffix(strings.ToLower(title), " blocks"), " ", "_")
				match := regexp.MustCompile(`block objects contain the following information within the (\w+) property`).FindStringSubmatch(desc)
				if len(match) != 0 {
					tagName = match[1]
				}
				builder.AddClass(getName(title), desc).AddDocProps(props...)
				builder.GetClass("Block").AddConfiguration(tagName, getName(title), desc)
			} else {
				panic(title)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := builder.Build("../types.block.go"); err != nil {
		t.Fatal(err)
	}
}

func TestPropertyItem(t *testing.T) {
	url := "https://developers.notion.com/reference/property-item-object"
	builder := NewBuilder().Add(CommentWithBreak(url))

	descRegex := regexp.MustCompile(`contain (.+) within the (\w+) property`)

	{
		code := jen.Type().Id("PropertyItem").Interface(jen.Id("GetType").Params().String()).Line()
		builder.Add(RawCoder{Value: code})
	}

	cases := []jen.Code{}

	err := Parse(url, func(title, desc string, props []DocProp) error {
		if title == "All property items" {
			builder.AddClass("PropertyItemCommon", desc).AddDocProps(props...)
			code := jen.Func().Params(jen.Id("i").Id("PropertyItemCommon")).Id("GetType").Params().String().Block(
				jen.Return().Id("i").Dot("Type"),
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

			cases = append(cases, jen.Case(jen.Lit(typeName)).Return().Op("&").Id(name).Block())

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
						prop.Type = jen.Index().Id("User")
					case "an array of file references":
						prop.Type = jen.Index().Id("File")
					case "an array of property_item objects":
						prop.Type = jen.Index().Id("DatePropertyItem")
					case "an array of multi-select option values":
						prop.Type = jen.Index().Id("MultiSelectOptionValues")
					case "an array of relation property items with a pagereferences":
						prop.Type = jen.Index().Id("RelationPropertyItem")
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
		t.Fatal(err)
	}

	// type PropertyItems []PropertyItem

	// func (items PropertyItems) UnmarshalJSON(data []byte) error {
	// 	fmt.Println(string(data))
	// 	return nil
	// }

	cases = append(cases, jen.Default().Panic(jen.Id("typeName")))

	code := jen.Func().Id("createPropertyItem").Params(jen.Id("typeName").String()).Id("PropertyItem").Block(
		jen.Switch(jen.Id("typeName")).Block(cases...),
	)
	builder.Add(RawCoder{Value: code})

	if err := builder.Build("../types.propertyitem.go"); err != nil {
		t.Fatal(err)
	}
}
