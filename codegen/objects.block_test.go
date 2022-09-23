package codegen

import (
	"regexp"
	"strings"
	"testing"

	"github.com/dave/jennifer/jen"
)

func TestBlockObject(t *testing.T) {
	t.Parallel()

	url := "https://developers.notion.com/reference/block"
	builder := NewBuilder().Add(CommentWithBreak(url))

	sections, err := ParseObjectDoc(url)
	if err != nil {
		t.Fatal(err)
	}

	descRegex := regexp.MustCompile(`block objects contain the following information within the (\w+) property`)

	for _, section := range sections {
		if heading := section.Heading; heading != nil {
			title := heading.Text
			desc := section.ParagraphText()

			switch title {
			case "Block object keys":
				obj := builder.AddClass("Block", desc)
				for _, param := range section.Parameters() {
					if param.Name == "{type}" {
						continue // 無視
					}
					if prop, err := param.Property(); err != nil {
						t.Fatal(err)
					} else {
						obj.AddField(prop)
					}
				}
			case "Block Type Object": // 無視
			case "Column List and Column Blocks":
				builder.AddClass("ColumnListBlocks", desc)
				builder.AddClass("ColumnBlocks", desc)
				builder.GetClass("Block").AddConfiguration("column_list", "ColumnListBlocks", desc)
				builder.GetClass("Block").AddConfiguration("column", "ColumnBlocks", desc)
			case "Synced Block blocks":
				for _, elem := range section.Elements {
					switch elem := elem.(type) {
					case *ParagraphElement:
						desc = elem.Content

					case *BlockParametersElement:
						topParam := (*elem)[0]
						if topParam.Name == "synced_from" && topParam.Type == "null" {
						} else if topParam.Name == "synced_from" && topParam.Type != "null" {
							obj := builder.AddClass("SyncedBlockBlocks", desc)
							builder.GetClass("Block").AddConfiguration("synced_block", "SyncedBlockBlocks", desc)
							for _, param := range *elem {
								if prop, err := param.Property(); err != nil {
									t.Fatal(err)
								} else {
									obj.AddField(prop)
								}
							}
						} else if topParam.Name == "type" {
							obj := builder.AddClass("SyncedFrom", desc)
							for _, param := range *elem {
								if prop, err := param.Property(); err != nil {
									t.Fatal(err)
								} else {
									obj.AddField(prop)
								}
							}
						} else {
							t.Fatal(topParam)
						}
					}
				}

			case "Image blocks", "Video blocks":
				for _, param := range section.Parameters() {
					if prop, err := param.Property(); err != nil {
						t.Fatal(err)
					} else {
						prop.Type = jen.Op("*").Id("ImageFile")
						prop.TypeSpecific = true
						prop.Description = desc
						builder.GetClass("Block").AddField(prop)
					}
				}
			default:
				if !strings.HasSuffix(title, " blocks") {
					t.Fatal(heading.Text)
				}

				tagName := strings.ReplaceAll(strings.TrimSuffix(strings.ToLower(title), " blocks"), " ", "_")
				if match := descRegex.FindStringSubmatch(desc); len(match) != 0 {
					tagName = match[1]
				}

				prop := Property{Name: tagName, Description: desc, TypeSpecific: true}
				if strings.HasPrefix(title, "Heading ") {
					obj := builder.GetClass("HeadingBlockData")
					if obj == nil {
						obj = builder.AddClass("HeadingBlockData", desc)
						for _, param := range section.Parameters() {
							if prop, err := param.Property(); err != nil {
								t.Fatal(err)
							} else {
								obj.AddField(prop)
							}
						}
					} else {
						obj.Comment += "\n" + desc
					}
					prop.Type = jen.Id(obj.Name)
				} else {
					if strings.Contains(desc, "do not contain any information within") {
						prop.Type = jen.Struct()
					} else {
						obj := builder.AddClass(getName(strings.TrimSuffix(title, "s")+" Data"), desc)
						prop.Type = jen.Id(obj.Name)
						for _, param := range section.Parameters() {
							if prop, err := param.Property(); err != nil {
								t.Fatal(err)
							} else {
								// childrenはomitemptyされることをAPI挙動で確認
								if param.Name == "children" {
									prop.OmitEmpty = true
								}
								obj.AddField(prop)
							}
						}
					}
				}
				builder.GetClass("Block").AddField(prop)
			}
		}
	}

	if err := builder.Build("../types.block.go"); err != nil {
		t.Fatal(err)
	}
}
