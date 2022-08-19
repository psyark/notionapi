package codegen

func BuildPage() error {
	url := "https://developers.notion.com/reference/page"
	builder := NewBuilder().Add(CommentWithBreak(url))

	err := Parse(url, func(title, desc string, props []DocProp) error {
		builder.AddClass("Page", desc).AddDocProps(props...)
		return nil
	})
	if err != nil {
		return err
	}

	return builder.Build("../types.page.go")
}
