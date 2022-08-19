package codegen

import (
	"regexp"
	"strings"
)

func BuildBlock() error {
	url := "https://developers.notion.com/reference/block"
	builder := NewBuilder().Add(CommentWithBreak(url))

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
				match := regexp.MustCompile(`block objects contain the following information within the (\w+) property`).FindStringSubmatch(desc)
				if len(match) != 0 {
					tagName = match[1]
				}
				builder.AddClass(getName(title), desc).AddDocProps(props...)
				builder.GetClass("Block").AddConfiguration(tagName, getName(title), desc)
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
