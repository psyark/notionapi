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
				prefix := strings.TrimSuffix(title, " objects")
				obj := builder.AddClass(prefix, desc)
				if err := obj.AddParams(nil, section.Parameters()...); err != nil {
					t.Fatal(err)
				}

				if title == "Link objects" {
					obj.AddField(&Property{Name: "url", Type: jen.String()})
				}

				builder.GetClass("RichText").AddConfiguration(nf_snake_case.String(prefix), prefix, desc)
			case strings.HasSuffix(title, " mentions"):
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
	}

	if err := builder.Build("../types.richtext.go"); err != nil {
		t.Fatal(err)
	}
}
