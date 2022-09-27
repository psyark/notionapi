package codegen

import (
	"testing"

	"github.com/dave/jennifer/jen"
)

func TestParentObject(t *testing.T) {
	t.Parallel()

	url := "https://developers.notion.com/reference/parent-object"
	builder := NewBuilder().Add(CommentWithBreak(url))

	sections, err := ParseObjectDoc(url)
	if err != nil {
		t.Fatal(err)
	}

	for _, section := range sections {
		heading := section.Heading
		if heading == nil {
			desc := section.AllParagraphText()
			builder.AddClass("Parent", desc).AddField(&Property{Name: "type", Type: jen.String()})
		} else {
			for _, param := range section.Parameters() {
				if param.Name == "type" {
					continue
				}
				if err := builder.GetClass("Parent").AddParams(&PropertyOption{TypeSpecific: true}, param); err != nil {
					t.Fatal(err)
				}
			}
		}
	}

	if err := builder.Build("../types.parent.go"); err != nil {
		t.Fatal(err)
	}
}
