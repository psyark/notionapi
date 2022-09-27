package codegen

import (
	"testing"

	"github.com/dave/jennifer/jen"
)

func TestPageObject(t *testing.T) {
	t.Parallel()

	url := "https://developers.notion.com/reference/page"
	builder := NewBuilder().Add(CommentWithBreak(url))

	sections, err := ParseObjectDoc(url)
	if err != nil {
		t.Fatal(err)
	}

	for _, section := range sections {
		heading := section.Heading
		if heading == nil {
			desc := section.ParagraphText()
			builder.AddClass("Page", desc)
		} else {
			for _, param := range section.Parameters() {
				if param.Name == "properties" {
					prop := &Property{
						Name:        "properties",
						Type:        jen.Map(jen.String()).Id("PropertyValue"),
						Description: param.Description,
					}
					builder.GetClass("Page").AddField(prop)
				} else {
					if err := builder.GetClass("Page").AddParams(nil, param); err != nil {
						t.Fatal(err)
					}
				}
			}
		}
	}

	if err := builder.Build("../types.page.go"); err != nil {
		t.Fatal(err)
	}
}
