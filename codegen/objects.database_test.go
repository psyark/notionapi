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
			desc := section.ParagraphText()
			builder.AddClass("Database", desc)
		} else {
			for _, param := range section.Parameters() {
				if param.Name == "properties*" {
					prop := Property{
						Name:        "properties",
						Type:        jen.Map(jen.String()).Id("Property"),
						Description: param.Description,
					}
					builder.GetClass("Database").AddField(prop)
				} else {
					prop, err := param.Property()
					if err != nil {
						t.Fatal(err)
					}
					builder.GetClass("Database").AddField(prop)
				}
			}
		}
	}

	if err := builder.Build("../types.database.go"); err != nil {
		t.Fatal(err)
	}
}
