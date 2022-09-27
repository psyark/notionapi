package codegen

import (
	"regexp"
	"strings"
	"testing"

	"github.com/dave/jennifer/jen"
)

func TestPropertyItemObject(t *testing.T) {
	t.Parallel()

	url := "https://developers.notion.com/reference/property-item-object"
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
			builder.AddClass("PropertyItem", desc)
			builder.AddClass("PaginatedPropertyItem", "")
			continue
		}

		title := heading.Text
		switch title {
		case "All property items":
			obj1 := builder.GetClass("PropertyItem").AddField(Comment(desc))
			obj2 := builder.GetClass("PaginatedPropertyItem")
			for _, param := range section.Parameters() {
				if param.Name != "next_url" {
					obj1.AddParams(nil, param)
				}
				if param.Name != "object" {
					obj2.AddParams(nil, param)
				}
			}
		case "Paginated property values":
			obj := builder.GetClass("PaginatedPropertyItem")
			obj.Comment = desc
			for _, name := range []string{"title", "rich_text", "relation", "rollup", "people"} {
				prop := &Property{Name: name, Type: jen.Struct(), TypeSpecific: true}
				if name == "rollup" {
					prop.Type = jen.Op("*").Id("RollupPropertyItemData")
				}
				obj.AddField(prop)
			}
		case "Select property values":
			builder.GetClass("PropertyItem").AddConfiguration("select", "SelectOption", desc)
			builder.GetClass("PropertyItem").AddConfiguration("status", "StatusOption", "undocumented")
		case "Multi-select option values": // ignore
		case "Date property values":
			builder.GetClass("PropertyItem").AddConfiguration("date", "DateValue", desc)
		case "Rollup property values", "Formula property values":
			prefix := strings.TrimSuffix(title, " property values")
			name := getName(prefix + " property item data")
			builder.GetClass("PropertyItem").AddConfiguration(strings.ToLower(prefix), name, desc)
			builder.AddClass(name, desc).AddParams(nil, section.Parameters()...)
		case "Incomplete rollup property values":
			builder.GetClass("RollupPropertyItemData").AddField(&Property{
				Name:         "incomplete",
				Type:         jen.Op("*").Struct(),
				Description:  desc,
				TypeSpecific: true,
			})
		default:
			if strings.HasSuffix(title, " formula property values") || strings.HasSuffix(title, " rollup property values") {
				if title == "Array rollup property values" {
					// ドキュメントに2箇所間違い
					prop := &Property{
						Name:         "array",
						Type:         jen.Index().Struct(),
						Description:  desc,
						TypeSpecific: true,
					}
					builder.GetClass("RollupPropertyItemData").AddField(prop)
					break
				}

				match := descRegex.FindStringSubmatch(desc)
				param := ObjectDocParameter{Name: match[2], Type: match[1], Description: desc}
				prop, err := param.Property(&PropertyOption{TypeSpecific: true})
				if err != nil {
					t.Fatal(err)
				}

				if strings.HasSuffix(title, " formula property values") {
					builder.GetClass("FormulaPropertyItemData").AddField(prop)
				} else if strings.HasSuffix(title, " rollup property values") {
					builder.GetClass("RollupPropertyItemData").AddField(prop)
				}
			} else if strings.HasSuffix(title, " property values") {
				name := getName(strings.Replace(title, " property values", " property item", 1))

				match := descRegex.FindStringSubmatch(desc)
				if len(match) == 0 {
					t.Fatal(desc)
				}

				if match[1] == "the following data" {
					builder.GetClass("PropertyItem").AddConfiguration(match[2], name+"Data", desc)
					builder.AddClass(name+"Data", "").AddParams(nil, section.Parameters()...)
				} else {
					param := ObjectDocParameter{Name: match[2], Type: match[1], Description: desc}

					switch param.Name {
					case "title", "rich_text", "relation", "people": // ドキュメントの "array of" は間違い
						prop := &Property{Name: param.Name, Description: param.Description, TypeSpecific: true}
						switch param.Type {
						case "array of user objects":
							prop.Type = jen.Op("*").Id("User")
						case "array of rich text objects":
							prop.Type = jen.Op("*").Id("RichText")
						case "array of relation property items with page references":
							prop.Type = jen.Op("*").Id("PageReference")
						default:
							t.Fatal(param)
						}
						builder.GetClass("PropertyItem").AddField(prop)
					case "multi_select":
						prop := &Property{Name: param.Name, Type: jen.Index().Id("SelectOption"), Description: param.Description, TypeSpecific: true}
						builder.GetClass("PropertyItem").AddField(prop)
					default:
						opt := &PropertyOption{TypeSpecific: true}
						if err := builder.GetClass("PropertyItem").AddParams(opt, param); err != nil {
							t.Fatal(err)
						}
					}
				}
			} else {
				t.Error(title)
			}
		}
	}

	if err := builder.Build("../types.propertyitem.go"); err != nil {
		t.Fatal(err)
	}
}
