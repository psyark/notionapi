package codegen

import (
	"fmt"
	"strings"

	"github.com/dave/jennifer/jen"
)

func BuildUser() error {
	url := "https://developers.notion.com/reference/user"
	builder := NewBuilder().Add(CommentWithBreak(url))

	err := Parse(url, func(title, desc string, props []DocProp) error {
		switch title {
		case "Where user objects appear in the API":
			builder.AddClass("User", desc)
		case "All users":
			builder.GetClass("User").AddField(Comment(desc)).AddDocProps(props...)
		case "People":
			person := builder.AddClass("Person", desc)
			for _, p := range props {
				if p.Name == "person" {
					builder.GetClass("User").AddField(Comment(p.Description)).AddConfiguration(p.Name, "Person", "")
				} else if strings.HasPrefix(p.Name, "person.") {
					p.Name = strings.TrimPrefix(p.Name, "person.")
					person.AddDocProps(p)
				} else {
					return fmt.Errorf("unknown property: %v", p.Name)
				}
			}
		case "Bots":
			if props[1].Name == "" && props[2].Name == "owner.type" && props[4].Name == "owner.type" {
				props = append(props[0:1], props[3:]...)
			} else {
				panic("ドキュメントが訂正されたようなので修正ロジックを見直してください")
			}

			for _, dp := range props {
				if dp.Name == "bot" {
					builder.GetClass("User").AddField(Comment(dp.Description)).AddConfiguration(dp.Name, "Bot", "")
					builder.AddClass("Bot", desc)
				} else if dp.Name == "owner" {
					builder.GetClass("Bot").AddField(&Property{
						Name:        dp.Name,
						Type:        jen.Op("*").Id("Owner"),
						Description: dp.Description,
					})
					builder.AddClass("Owner", desc)
				} else if strings.HasPrefix(dp.Name, "owner.") {
					dp.Name = strings.TrimPrefix(dp.Name, "owner.")
					p := dp.Property()
					builder.GetClass("Owner").AddField(p)
				} else {
					return fmt.Errorf("unknown property: %v", dp.Name)
				}
			}
		default:
			return fmt.Errorf("unknown title: %v", title)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return builder.Build("../types.user.go")
}
