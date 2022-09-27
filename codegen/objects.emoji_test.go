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
		desc := section.AllParagraphText()

		if heading == nil {
			builder.AddClass("Emoji", desc).Implement("fileOrEmoji")
		} else {
			if err := builder.GetClass("Emoji").AddParams(nil, section.Parameters()...); err != nil {
				t.Fatal(err)
			}
		}
	}

	if err := builder.Build("../types.emoji.go"); err != nil {
		t.Fatal(err)
	}
}
