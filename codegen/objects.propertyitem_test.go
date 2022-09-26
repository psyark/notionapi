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
		desc := section.ParagraphText()

		if heading == nil {
			builder.AddClass("PropertyItem", desc)
			continue
		}

		title := heading.Text
		switch title {
		case "All property items":
			obj := builder.GetClass("PropertyItem").AddField(Comment(desc))
			for _, param := range section.Parameters() {
				if param.Name == "next_url" {
					prop, err := param.Property()
					if err != nil {
						t.Fatal(err)
					}

					prop.OmitEmpty = true
					obj.AddField(prop)
				} else {
					obj.AddParams(param)
				}
			}
		case "Paginated property values": // ignore
		case "Multi-select option values": // ignore
		case "Rollup property values", "Formula property values":
			prefix := strings.TrimSuffix(title, " property values")
			name := getName(prefix + " property item data")
			builder.GetClass("PropertyItem").AddConfiguration2(strings.ToLower(prefix), name, desc)
			builder.AddClass(name, desc).AddParams(section.Parameters()...)
		case "Incomplete rollup property values":
			builder.GetClass("RollupPropertyItemData").AddField(Property{
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
				prop, err := param.Property()
				if err != nil {
					t.Fatal(err)
				}

				prop.TypeSpecific = true
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
					builder.GetClass("PropertyItem").AddConfiguration2(match[2], name+"Data", desc)
					builder.AddClass(name+"Data", "").AddParams(section.Parameters()...)
					if match[2] == "select" { // ドキュメントに抜けているstatusを作る
						name = strings.Replace(name, "Select", "Status", 1)
						builder.GetClass("PropertyItem").AddConfiguration2("status", name+"Data", desc)
						builder.AddClass(name+"Data", "").AddParams(section.Parameters()...)
					}
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
					default:
						prop, err := param.Property()
						if err != nil {
							t.Fatal(err)
						}

						prop.TypeSpecific = true
						builder.GetClass("PropertyItem").AddField(prop)
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
