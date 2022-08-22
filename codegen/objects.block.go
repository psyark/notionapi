package codegen

import (
	"regexp"
	"strings"

	"github.com/dave/jennifer/jen"
)

func BuildBlock() error {
	url := "https://developers.notion.com/reference/block"
	builder := NewBuilder().Add(CommentWithBreak(url))

	descRegex := regexp.MustCompile(`block objects contain the following information within the (\w+) property`)

	err := Parse(url, func(title, desc string, props []DocProp) error {
		switch title {
		case "Block object keys":
			object := builder.AddClass("Block", desc)
			for _, p := range props {
				if p.Name != "{type}" {
					object.AddDocProps(p)
				}
			}
		case "Block Type Object":
		case "Column List and Column Blocks":
			builder.AddClass("ColumnListBlocks", desc)
			builder.AddClass("ColumnBlocks", desc)
			builder.GetClass("Block").AddConfiguration("column_list", "ColumnListBlocks", desc)
			builder.GetClass("Block").AddConfiguration("column", "ColumnBlocks", desc)
		case "Synced Block blocks":
			builder.AddClass("SyncedBlockBlocks", desc).AddDocProps(props[2])
			builder.GetClass("Block").AddConfiguration("synced_block", "SyncedBlockBlocks", desc)
			builder.AddClass("SyncedFrom", desc).AddDocProps(props[3:]...)
		default:
			if strings.HasSuffix(title, " Blocks") || strings.HasSuffix(title, " blocks") {
				tagName := strings.ReplaceAll(strings.TrimSuffix(strings.ToLower(title), " blocks"), " ", "_")
				if match := descRegex.FindStringSubmatch(desc); len(match) != 0 {
					tagName = match[1]
				}

				var object *Class
				if strings.HasPrefix(title, "Heading ") {
					object = builder.GetClass("HeadingBlockData")
					if object == nil {
						object = builder.AddClass("HeadingBlockData", desc).AddDocProps(props...)
						object.AddField(Property{Name: "is_toggleable", Type: jen.Bool(), Description: "undocumented"})
					} else {
						object.Comment += "\n" + desc
					}
				} else {

					object = builder.AddClass(getName(strings.TrimSuffix(title, "s")+" Data"), desc).AddDocProps(props...)
				}

				builder.GetClass("Block").AddConfiguration(tagName, object.Name, desc)
			} else {
				panic(title)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return builder.Build("../types.block.go")
}
