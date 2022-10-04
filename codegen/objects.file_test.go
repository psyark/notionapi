package codegen

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/dave/jennifer/jen"
)

func TestFileObject(t *testing.T) {
	t.Parallel()

	url := "https://developers.notion.com/reference/file-object"
	builder := NewBuilder().Add(CommentWithBreak(url))

	sections, err := ParseObjectDoc(url)
	if err != nil {
		t.Fatal(err)
	}

	file := builder.AddClass("File", sections[0].AllParagraphText()).Implement("fileOrEmoji")

	for _, section := range sections[1:] {
		title := section.Heading.Text
		desc := section.AllParagraphText()

		switch title {
		case "All file objects":
		ELEMENTS_LOOP:
			for _, elem := range section.Elements {
				switch elem := elem.(type) {
				case *ParagraphElement:
					file.AddField(Comment(elem.Content))
				case *BlockParametersElement:
					for _, param := range *elem {
						if param.Name == "{type}" {
							param.Name = "type"
						}
						if err := file.AddParams(nil, param); err != nil {
							t.Fatal(err)
						}
					}

					// nilの場合はomitし、[]の場合はomitしないため *[]RichText とする
					caption := &Property{Name: "caption", Type: jen.Op("*").Id("RichTextArray"), Description: "undocumented", OmitEmpty: true}
					file.AddField(caption).AddLine()
				case *BlockAPIHeaderElement:
					if elem.Title == "Externally hosted files vs. Files hosted by Notion" {
						break ELEMENTS_LOOP
					}
				default:
					fmt.Println(reflect.TypeOf(elem))
				}
			}

		case "Files uploaded to Notion objects":
			err := builder.AddClass("FilesUploadedToNotionData", desc).AddParams(nil, section.Parameters()...)
			if err != nil {
				t.Fatal(err)
			}
			file.AddConfiguration("file", "FilesUploadedToNotionData", desc)
		case "External file objects":
			err := builder.AddClass("ExternalFileData", desc).AddParams(nil, section.Parameters()...)
			if err != nil {
				t.Fatal(err)
			}
			file.AddConfiguration("external", "ExternalFileData", desc)
		default:
			t.Error(title)
		}
	}

	if err := builder.Build("../types.file.go"); err != nil {
		t.Fatal(err)
	}
}
