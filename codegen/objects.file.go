package codegen

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dave/jennifer/jen"
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
				className := getName(strings.TrimSuffix(title, " objects")) + "Data"
				prop := Property{Name: match[1], Type: jen.Id(className), Description: desc, TypeSpecific: true}
				builder.AddClass(className, title).AddDocProps(props...)
				builder.GetClass("File").AddField(prop)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return builder.Build("../types.file.go")
}
