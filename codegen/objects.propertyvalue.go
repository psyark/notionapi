package codegen

func BuildPropertyValue() error {
	url := "https://developers.notion.com/reference/property-value-object"
	builder := NewBuilder().Add(CommentWithBreak(url))

	err := Parse(url, func(title, desc string, props []DocProp) error {
		if title == "All property values" {
			builder.AddClass("PropertyValue", desc).AddDocProps(props...)
		} else {
			builder.GetClass("PropertyValue").AddField(Comment(title + ": " + desc)) // TODO
		}
		return nil
	})
	if err != nil {
		return err
	}

	return builder.Build("../types.propertyvalue.go")
}
