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

	page := builder.AddClass("Page", sections[0].AllParagraphText())

	for _, section := range sections[1:] {
		for _, param := range section.Parameters() {
			if param.Name == "properties" {
				prop := &Property{
					Name:        "properties",
					Type:        jen.Map(jen.String()).Id("PropertyValue"),
					Description: param.Description,
				}
				page.AddField(prop)
			} else {
				if err := page.AddParams(nil, param); err != nil {
					t.Fatal(err)
				}
			}
		}
	}

	if err := builder.Build("../types.page.go"); err != nil {
		t.Fatal(err)
	}
}
