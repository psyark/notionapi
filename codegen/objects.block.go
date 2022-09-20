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
			for _, dp := range props {
				if dp.Name == "{type}" {
					continue
				}
				prop := dp.Property()
				if dp.Type == "boolean" {
					prop.Type = jen.Op("*").Bool()
				}
				prop.OmitEmpty = true
				object.AddField(prop)
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
		case "Image blocks", "Video blocks":
			prop := props[0].Property()
			prop.TypeSpecific = true
			prop.Description = desc
			prop.Type = jen.Id("File")
			builder.GetClass("Block").AddField(prop)
		default:
			if strings.HasSuffix(title, " Blocks") || strings.HasSuffix(title, " blocks") {
				tagName := strings.ReplaceAll(strings.TrimSuffix(strings.ToLower(title), " blocks"), " ", "_")
				if match := descRegex.FindStringSubmatch(desc); len(match) != 0 {
					tagName = match[1]
				}

				prop := Property{Name: tagName, Description: desc, TypeSpecific: true}
				if strings.HasPrefix(title, "Heading ") {
					object := builder.GetClass("HeadingBlockData")
					if object == nil {
						object = builder.AddClass("HeadingBlockData", desc).AddDocProps(props...)
					} else {
						object.Comment += "\n" + desc
					}
					prop.Type = jen.Id(object.Name)
				} else {
					if strings.Contains(desc, "do not contain any information within") {
						prop.Type = jen.Struct()
					} else {
						object := builder.AddClass(getName(strings.TrimSuffix(title, "s")+" Data"), desc).AddDocProps(props...)
						prop.Type = jen.Id(object.Name)
					}
				}
				builder.GetClass("Block").AddField(prop)

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
