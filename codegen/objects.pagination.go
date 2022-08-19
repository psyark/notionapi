package codegen

import (
	"regexp"

	"github.com/dave/jennifer/jen"
)

func BuildPagination() error {
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
		return err
	}

	return builder.Build("../types.pagination.go")
}
