package codegen

import (
	"fmt"
	"strings"

	"github.com/dave/jennifer/jen"
)

func BuildRichText() error {
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

			prop := Property{Name: tagName, Type: jen.Id(className), Description: desc, TypeSpecific: true}
			builder.GetClass("RichText").AddField(prop)
		default:
			return fmt.Errorf("unknown title: %v", title)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return builder.Build("../types.richtext.go")
}
