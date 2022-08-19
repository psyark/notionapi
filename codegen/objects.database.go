package codegen

func BuildDatabase() error {
	url := "https://developers.notion.com/reference/database"
	builder := NewBuilder().Add(CommentWithBreak(url))

	err := Parse(url, func(title, desc string, props []DocProp) error {
		builder.AddClass("Database", desc).AddDocProps(props...)
		return nil
	})
	if err != nil {
		return err
	}

	return builder.Build("../types.database.go")
}
