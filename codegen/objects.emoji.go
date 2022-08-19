package codegen

func BuildEmoji() error {
	url := "https://developers.notion.com/reference/emoji-object"
	builder := NewBuilder().Add(CommentWithBreak(url))

	err := Parse(url, func(title, desc string, props []DocProp) error {
		builder.AddClass("Emoji", desc).AddDocProps(props...)
		return nil
	})
	if err != nil {
		return err
	}

	return builder.Build("../types.emoji.go")
}
