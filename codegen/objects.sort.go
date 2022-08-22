package codegen

func BuildSort() error {
	url := "https://developers.notion.com/reference/post-database-query-sort"
	builder := NewBuilder().Add(CommentWithBreak(url))

	err := Parse(url, func(title, desc string, props []DocProp) error {
		builder.Add(Comment(title))
		builder.Add(Comment(desc))
		// builder.AddClass("Sort", desc).AddDocProps(props...)
		return nil
	})
	if err != nil {
		return err
	}

	return builder.Build("../types.sort.go")
}
