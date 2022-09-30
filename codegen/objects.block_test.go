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

	descRegex := regexp.MustCompile(`block objects contain the following information within the (\w+) property`)

	sections, err := ParseObjectDoc(url)
	if err != nil {
		t.Fatal(err)
	}

	block := builder.AddClass("Block", sections[0].AllParagraphText())

	for _, section := range sections[1:] {
		heading := section.Heading
		desc := section.AllParagraphText()

		title := heading.Text

		switch title {
		case "Block object keys":
			obj := block.AddField(Comment(desc))
			for _, param := range section.Parameters() {
				if param.Name == "{type}" {
					continue // 無視
				}

				opt := &PropertyOption{OmitEmpty: true, Nullable: param.Type == "boolean"}
				if prop, err := param.Property(opt); err != nil {
					t.Fatal(err)
				} else {
					obj.AddField(prop)
				}
			}
			obj.AddLine()

		case "Block Type Object": // 無視
		case "Column List and Column Blocks":
			builder.AddClass("ColumnListBlocks", desc)
			builder.AddClass("ColumnBlocks", desc)
			block.AddConfiguration("column_list", "ColumnListBlocks", desc)
			block.AddConfiguration("column", "ColumnBlocks", desc)
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
						block.AddConfiguration("synced_block", "SyncedBlockBlocks", desc)
						if err := obj.AddParams(nil, *elem...); err != nil {
							t.Fatal(err)
						}
					} else if topParam.Name == "type" {
						obj := builder.AddClass("SyncedFrom", desc)
						if err := obj.AddParams(nil, *elem...); err != nil {
							t.Fatal(err)
						}
					} else {
						t.Fatal(topParam)
					}
				}
			}

		case "Image blocks", "Video blocks":
			for _, param := range section.Parameters() {
				if prop, err := param.Property(&PropertyOption{TypeSpecific: true}); err != nil {
					t.Fatal(err)
				} else {
					prop.Description = desc
					block.AddField(prop)
				}
			}
		default:
			if !strings.HasSuffix(title, " blocks") {
				t.Fatal(heading.Text)
			}

			tagName := nf_snake_case.String(strings.TrimSuffix(title, " blocks"))
			if match := descRegex.FindStringSubmatch(desc); len(match) != 0 {
				tagName = match[1]
			}

			prop := &Property{Name: tagName, Description: desc, TypeSpecific: true}
			if strings.HasPrefix(title, "Heading ") {
				obj := builder.GetClass("HeadingBlockData")
				if obj == nil {
					obj = builder.AddClass("HeadingBlockData", desc)
					if err := obj.AddParams(nil, section.Parameters()...); err != nil {
						t.Fatal(err)
					}
				} else {
					obj.Comment += "\n" + desc
				}
				prop.Type = jen.Id(obj.Name)
			} else if strings.Contains(desc, "do not contain any information within") {
				prop.Type = jen.Struct()
			} else {
				obj := builder.AddClass(nfCamelCase.String(strings.TrimSuffix(title, "s"))+"Data", desc)
				prop.Type = jen.Id(obj.Name)
				for _, param := range section.Parameters() {
					opt := &PropertyOption{OmitEmpty: param.Name == "children"} // childrenはomitemptyされることをAPI挙動で確認
					if err := obj.AddParams(opt, param); err != nil {
						t.Fatal(err)
					}
				}
				if title == "Embed blocks" {
					obj.AddField(&Property{Name: "caption", Type: jen.Index().Id("RichText"), Description: "undocumented"})
				}
			}

			block.AddField(prop)
		}
	}

	if err := builder.Build("../types.block.go"); err != nil {
		t.Fatal(err)
	}
}
