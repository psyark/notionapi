package codegen

import (
	"github.com/dave/jennifer/jen"
)

func BuildParent() error {
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
		return err
	}

	return builder.Build("../types.parent.go")
}
