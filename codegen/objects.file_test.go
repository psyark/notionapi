package codegen

import (
	"testing"
)

func TestFileObject(t *testing.T) {
	t.Parallel()

	url := "https://developers.notion.com/reference/file-object"
	builder := NewBuilder().Add(CommentWithBreak(url))

	sections, err := ParseObjectDoc(url)
	if err != nil {
		t.Fatal(err)
	}

	for _, section := range sections {
		heading := section.Heading
		desc := section.AllParagraphText()

		if heading == nil {
			builder.AddClass("File", desc).Implement("fileOrEmoji")
		} else {
			title := heading.Text

			switch title {
			case "All file objects":
				for _, param := range section.Parameters() {
					if param.Name == "" {
						continue
					} else if param.Name == "{type}" {
						param.Name = "type"
					}
					if err := builder.GetClass("File").AddParams(nil, param); err != nil {
						t.Fatal(err)
					}
				}
			case "Files uploaded to Notion objects":
				err := builder.AddClass("FilesUploadedToNotionData", desc).AddParams(nil, section.Parameters()...)
				if err != nil {
					t.Fatal(err)
				}
				builder.GetClass("File").AddConfiguration("file", "FilesUploadedToNotionData", desc)
			case "External file objects":
				err := builder.AddClass("ExternalFileData", desc).AddParams(nil, section.Parameters()...)
				if err != nil {
					t.Fatal(err)
				}
				builder.GetClass("File").AddConfiguration("external", "ExternalFileData", desc)
			default:
				t.Error(title)
			}
		}
	}

	if err := builder.Build("../types.file.go"); err != nil {
		t.Fatal(err)
	}
}
