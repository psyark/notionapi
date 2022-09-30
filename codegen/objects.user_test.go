package codegen

import (
	"strings"
	"testing"

	"github.com/dave/jennifer/jen"
)

func TestUserObject(t *testing.T) {
	t.Parallel()

	url := "https://developers.notion.com/reference/user"
	builder := NewBuilder().Add(CommentWithBreak(url))

	sections, err := ParseObjectDoc(url)
	if err != nil {
		t.Fatal(err)
	}

	for _, section := range sections {
		heading := section.Heading
		if heading == nil {
			continue
		}

		title := heading.Text
		desc := section.AllParagraphText()

		switch title {
		case "Where user objects appear in the API":
			builder.AddClass("User", desc)
			builder.AddClass("PartialUser", desc)
		case "All users":
			obj1 := builder.GetClass("User").AddField(Comment(desc))
			obj2 := builder.GetClass("PartialUser").AddField(Comment(desc))
			for _, param := range section.Parameters() {
				// ドキュメントの optional は間違いと思われる
				if param.Type == "string (optional, enum)" && param.Name == "type" {
					param.Type = "string (enum)"
				}

				prop, err := param.Property(nil)
				if err != nil {
					t.Fatal(err)
				}

				obj1.AddField(prop)
				if strings.HasSuffix(param.Name, "*") {
					obj2.AddField(prop)
				}
			}
		case "People":
			for _, param := range section.Parameters() {
				if param.Name == "person" {
					prop := &Property{
						Name:         param.Name,
						Type:         jen.Op("*").Id("Person"),
						TypeSpecific: true,
						Description:  param.Description,
					}
					builder.GetClass("User").AddLine().AddField(prop)
					builder.AddClass("Person", desc)
				} else if strings.HasPrefix(param.Name, "person.") {
					param.Name = strings.TrimPrefix(param.Name, "person.")
					if err := builder.GetClass("Person").AddParams(nil, param); err != nil {
						t.Fatal(err)
					}
				} else {
					t.Error(param.Name)
				}
			}
		case "Bots":
			for _, elem := range section.Elements {
				switch elem := elem.(type) {
				case *ParagraphElement:
					desc = elem.Content
				case *BlockParametersElement:
					topParam := (*elem)[0]
					switch topParam.Name {
					case "bot":
						prop := &Property{
							Name:         "bot",
							Type:         jen.Op("*").Id("Bot"),
							Description:  topParam.Description,
							TypeSpecific: true,
						}
						builder.GetClass("User").AddField(prop)
						builder.AddClass("Bot", desc)
					case "owner":
						builder.AddClass("Owner", desc)
						for _, param := range *elem {
							if param.Name == "owner" {
								prop := &Property{
									Name:        param.Name,
									Type:        jen.Op("*").Id("Owner"),
									Description: param.Description,
									OmitEmpty:   true, // APIで挙動確認済
								}
								builder.GetClass("Bot").AddField(prop)
							} else {
								param.Name = strings.TrimPrefix(param.Name, "owner.")
								if err := builder.GetClass("Owner").AddParams(nil, param); err != nil {
									t.Fatal(err)
								}
							}
						}
					}
				}
			}
		default:
			t.Error(title)
		}
	}

	if err := builder.Build("../types.user.go"); err != nil {
		t.Fatal(err)
	}
}
