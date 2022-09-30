package codegen

import (
	"testing"

	"github.com/dave/jennifer/jen"
)

func TestDatabaseObject(t *testing.T) {
	t.Parallel()

	url := "https://developers.notion.com/reference/database"
	builder := NewBuilder().Add(CommentWithBreak(url))

	sections, err := ParseObjectDoc(url)
	if err != nil {
		t.Fatal(err)
	}

	database := builder.AddClass("Database", sections[0].AllParagraphText())

	for _, section := range sections[1:] {
		for _, param := range section.Parameters() {
			if param.Name == "properties*" {
				prop := &Property{
					Name:        "properties",
					Type:        jen.Map(jen.String()).Id("Property"),
					Description: param.Description,
				}
				database.AddField(prop)
			} else {
				if err := database.AddParams(nil, param); err != nil {
					t.Fatal(err)
				}
			}
		}
	}

	if err := builder.Build("../types.database.go"); err != nil {
		t.Fatal(err)
	}
}
