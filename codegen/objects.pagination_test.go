package codegen

import (
	"regexp"
	"testing"

	"github.com/dave/jennifer/jen"
)

func TestPaginationObject(t *testing.T) {
	t.Parallel()

	url := "https://developers.notion.com/reference/pagination"
	builder := NewBuilder().Add(CommentWithBreak(url))

	sections, err := ParseObjectDoc(url)
	if err != nil {
		t.Fatal(err)
	}

	for _, section := range sections {
		heading := section.Heading
		if heading == nil {
		} else if heading.Text == "Responses" {
			desc := section.ParagraphText()
			obj := builder.AddClass("Pagination", desc)

			for _, param := range section.Parameters() {
				switch param.Name {
				case "results", "{type}": // ignore
				default:
					prop, err := param.Property()
					if err != nil {
						t.Fatal(err)
					}
					if param.Name == "next_cursor" {
						prop.Type = jen.Op("*").String() // next_cursor might be null.
					} else if param.Name == "type" {
						match := regexp.MustCompile(` "(\w+)"`).FindAllStringSubmatch(param.Description, -1)
						for _, m := range match {
							name := getName(m[1])
							pagi := builder.AddClass(name+"Pagination", "").AddField(
								AnonymousField("Pagination"),
								Property{Name: "results", Type: jen.Index().Id(name)},
							)
							if name == "PropertyItem" {
								pagi.AddField(Property{Name: m[1], Type: jen.Id("PaginatedPropertyItem")})
							} else {
								pagi.AddField(Property{Name: m[1], Type: jen.Struct()})
							}
						}
					}
					obj.AddField(prop)
				}
			}
		}
	}

	if err := builder.Build("../types.pagination.go"); err != nil {
		t.Fatal(err)
	}
}
