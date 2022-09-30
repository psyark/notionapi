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

	parent := builder.AddClass("Parent", sections[0].AllParagraphText()).AddField(&Property{Name: "type", Type: jen.String()})

	for _, section := range sections[1:] {
		for _, param := range section.Parameters() {
			if param.Name == "type" {
				continue
			}
			if err := parent.AddParams(&PropertyOption{TypeSpecific: true}, param); err != nil {
				t.Fatal(err)
			}
		}
	}

	if err := builder.Build("../types.parent.go"); err != nil {
		t.Fatal(err)
	}
}
