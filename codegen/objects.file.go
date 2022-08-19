package codegen

import (
	"fmt"
	"regexp"
	"strings"
)

func BuildFile() error {
	url := "https://developers.notion.com/reference/file-object"
	builder := NewBuilder().Add(CommentWithBreak(url))

	err := Parse(url, func(title, desc string, props []DocProp) error {
		switch title {
		case "All file objects":
			for i := range props {
				if props[i].Name == "{type}" {
					props[i].Name = "type"
				}
			}
			builder.AddClass("File", desc).AddDocProps(props...)
		case "Externally hosted files vs. Files hosted by Notion":
			builder.GetClass("File").Comment += fmt.Sprintf("\n\n%v\n%v", title, desc)
		default:
			match := regexp.MustCompile(`have a type of "(\w+)"`).FindStringSubmatch(desc)
			if len(match) != 0 {
				className := getName(strings.ReplaceAll(strings.TrimSuffix(title, " objects"), " ", "_"))
				builder.AddClass(className, title+"\n"+desc).AddDocProps(props...)
				builder.GetClass("File").AddConfiguration(match[1], className, desc)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return builder.Build("../types.file.go")
}
