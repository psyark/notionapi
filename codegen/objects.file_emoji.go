package codegen

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dave/jennifer/jen"
)

func BuildFileOrEmoji() error {
	url1 := "https://developers.notion.com/reference/file-object"
	url2 := "https://developers.notion.com/reference/emoji-object"
	builder := NewBuilder().Add(CommentWithBreak(url1))

	builder.AddClass("FileOrEmoji", "")

	err := Parse(url1, func(title, desc string, props []DocProp) error {
		switch title {
		case "All file objects":
			for i := range props {
				if props[i].Name == "{type}" {
					props[i].Name = "type"
				}
			}
			builder.AddClass("File", desc).AddDocProps(props...)
			builder.GetClass("FileOrEmoji").AddDocProps(props...)
		case "Externally hosted files vs. Files hosted by Notion":
			builder.GetClass("File").Comment += fmt.Sprintf("\n\n%v\n%v", title, desc)
		default:
			match := regexp.MustCompile(`have a type of "(\w+)"`).FindStringSubmatch(desc)
			if len(match) != 0 {
				className := getName(strings.TrimSuffix(title, " objects")) + "Data"
				prop := Property{Name: match[1], Type: jen.Id(className), Description: desc, TypeSpecific: true}
				builder.AddClass(className, title).AddDocProps(props...)
				builder.GetClass("File").AddField(prop)
				builder.GetClass("FileOrEmoji").AddField(prop)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	builder.Add(CommentWithBreak(url2))

	err = Parse(url2, func(title, desc string, props []DocProp) error {
		builder.AddClass("Emoji", desc).AddDocProps(props...)
		for _, dp := range props {
			if dp.Name != "type" {
				prop := dp.Property()
				prop.TypeSpecific = true
				builder.GetClass("FileOrEmoji").AddField(prop)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return builder.Build("../types.file_emoji.go")
}
