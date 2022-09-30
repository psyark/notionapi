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

	emoji := builder.AddClass("Emoji", sections[0].AllParagraphText()).Implement("fileOrEmoji")

	for _, section := range sections[1:] {
		if err := emoji.AddParams(nil, section.Parameters()...); err != nil {
			t.Fatal(err)
		}
	}

	if err := builder.Build("../types.emoji.go"); err != nil {
		t.Fatal(err)
	}
}
