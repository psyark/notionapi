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

	rich_text := builder.AddClass("RichText", sections[0].AllParagraphText())

	builder.Add(RawCoder{jen.Type().Id("RichTextArray").Index().Id("RichText")})
	builder.Add(RawCoder{jen.Func().Params(jen.Id("a").Id("RichTextArray")).Id("PlainText").Params().String().Block(
		jen.Id("text").Op(":=").Lit(""),
		jen.For(jen.List(jen.Id("_"), jen.Id("rt")).Op(":=").Range().Id("a").Block(
			jen.Id("text").Op("+=").Id("rt").Dot("PlainText"),
		)),
		jen.Return().Id("text"),
	)})

	for _, section := range sections[1:] {
		title := section.Heading.Text
		desc := section.AllParagraphText()

		switch title {
		case "All rich text":
			obj := rich_text.AddField(Comment(desc))
			for _, param := range section.Parameters() {
				prop, err := param.Property(&PropertyOption{OmitEmpty: param.Name == "annotations"})
				if err != nil {
					t.Fatal(err)
				}
				obj.AddField(prop)
			}

			obj.AddLine()
		case "Annotations":
			obj := builder.AddClass(title, desc)
			if err := obj.AddParams(nil, section.Parameters()...); err != nil {
				t.Fatal(err)
			}
		case "Text objects", "Link objects", "Mention objects", "Equation objects":
			prefix := strings.TrimSuffix(title, " objects")
			obj := builder.AddClass(prefix, desc)
			if err := obj.AddParams(nil, section.Parameters()...); err != nil {
				t.Fatal(err)
			}

			if title == "Link objects" {
				obj.AddField(&Property{Name: "url", Type: jen.String()})
			}

			rich_text.AddConfiguration(nf_snake_case.String(prefix), prefix, desc)
		case "User mentions", "Page mentions", "Database mentions", "Date mentions", "Template mentions", "Link Preview mentions":
			mention := builder.GetClass("Mention")

			switch title {
			case "User mentions":
				mention.AddField(&Property{Name: "user", Type: jen.Op("*").Id("User"), Description: desc, TypeSpecific: true})
			case "Page mentions":
				mention.AddField(&Property{Name: "page", Type: jen.Op("*").Id("PageReference"), Description: desc, TypeSpecific: true})
			case "Database mentions":
				mention.AddField(&Property{Name: "database", Type: jen.Op("*").Id("PageReference"), Description: desc, TypeSpecific: true})
			case "Date mentions":
				mention.AddField(&Property{Name: "date", Type: jen.Op("*").Id("DateValue"), Description: desc, TypeSpecific: true})
			case "Template mentions": // TODO
				mention.AddField(Comment("TODO: " + title))
			case "Link Preview mentions":
				mention.AddField(&Property{Name: "link_preview", Type: jen.Op("*").Id("LinkPreview"), Description: desc, TypeSpecific: true})
				dataObj := builder.AddClass("LinkPreview", desc)
				if err := dataObj.AddParams(nil, section.Parameters()...); err != nil {
					t.Fatal(err)
				}
			default:
				t.Error(desc)
			}

		default:
			t.Error(title)
		}
	}

	if err := builder.Build("../types.richtext.go"); err != nil {
		t.Fatal(err)
	}
}
