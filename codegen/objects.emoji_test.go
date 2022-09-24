package codegen

import (
	"testing"
)

func TestEmojiObject(t *testing.T) {
	t.Parallel()

	url := "https://developers.notion.com/reference/emoji-object"
	builder := NewBuilder().Add(CommentWithBreak(url))

	sections, err := ParseObjectDoc(url)
	if err != nil {
		t.Fatal(err)
	}

	for _, section := range sections {
		heading := section.Heading
		desc := section.ParagraphText()

		if heading == nil {
			builder.AddClass("Emoji", desc)
		} else {
			err := builder.GetClass("Emoji").AddParams(section.Parameters()...)
			if err != nil {
				t.Fatal(err)
			}
		}
	}

	if err := builder.Build("../types.emoji.go"); err != nil {
		t.Fatal(err)
	}
}
