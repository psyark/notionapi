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

	for _, section := range sections {
		heading := section.Heading
		if heading == nil {
			desc := section.AllParagraphText()
			builder.AddClass("Database", desc)
		} else {
			for _, param := range section.Parameters() {
				if param.Name == "properties*" {
					prop := &Property{
						Name:        "properties",
						Type:        jen.Map(jen.String()).Id("Property"),
						Description: param.Description,
					}
					builder.GetClass("Database").AddField(prop)
				} else {
					if err := builder.GetClass("Database").AddParams(nil, param); err != nil {
						t.Fatal(err)
					}
				}
			}
		}
	}

	if err := builder.Build("../types.database.go"); err != nil {
		t.Fatal(err)
	}
}
