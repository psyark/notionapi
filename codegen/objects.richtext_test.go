package codegen

import (
	"strings"
	"testing"

	"github.com/dave/jennifer/jen"
)

func TestRichTextObject(t *testing.T) {
	t.Parallel()

	url := "https://developers.notion.com/reference/rich-text"
	builder := NewBuilder().Add(CommentWithBreak(url))

	sections, err := ParseObjectDoc(url)
	if err != nil {
		t.Fatal(err)
	}

	for _, section := range sections {
		heading := section.Heading
		if heading == nil {
			desc := section.AllParagraphText()
			builder.AddClass("RichText", desc)
		} else {
			title := heading.Text
			desc := section.AllParagraphText()

			switch {
			case title == "All rich text":
				obj := builder.GetClass("RichText").AddField(Comment(desc))
				if err := obj.AddParams(nil, section.Parameters()...); err != nil {
					t.Fatal(err)
				}
				obj.AddField(RawCoder{jen.Line()})
			case title == "Annotations":
				obj := builder.AddClass(title, desc)
				if err := obj.AddParams(nil, section.Parameters()...); err != nil {
					t.Fatal(err)
				}
			case strings.HasSuffix(title, " objects"):
				obj := builder.AddClass(strings.TrimSuffix(title, " objects"), desc)
				if err := obj.AddParams(nil, section.Parameters()...); err != nil {
					t.Fatal(err)
				}

				if title == "Link objects" {
					obj.AddField(&Property{Name: "url", Type: jen.String()})
				}

				builder.GetClass("RichText").AddConfiguration(strings.ToLower(obj.Name), obj.Name, desc)
			case strings.HasSuffix(title, " mentions"):
				err := builder.GetClass("Mention").AddField(Comment(desc)).AddParams(nil, section.Parameters()...)
				if err != nil {
					t.Fatal(err)
				}
			default:
				t.Error(title)
			}
		}
	}

	if err := builder.Build("../types.richtext.go"); err != nil {
		t.Fatal(err)
	}
}
