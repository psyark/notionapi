package codegen

import (
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestMethods(t *testing.T) {
	type methodDef struct {
		DocURL  string
		Type    MethodCoderType
		Correct func(*MethodCoder)
	}

	methods := []methodDef{
		{DocURL: "https://developers.notion.com/reference/post-database-query", Type: &ReturnsStructRef{"PagePagination"}},
		{DocURL: "https://developers.notion.com/reference/create-a-database", Type: &ReturnsStructRef{"Database"}},
		{DocURL: "https://developers.notion.com/reference/update-a-database", Type: &ReturnsStructRef{"Database"}},
		{DocURL: "https://developers.notion.com/reference/retrieve-a-database", Type: &ReturnsStructRef{"Database"}},
		// https://developers.notion.com/reference/get-databases : deprecated
		{DocURL: "https://developers.notion.com/reference/retrieve-a-page", Type: &ReturnsStructRef{"Page"}},
		{DocURL: "https://developers.notion.com/reference/post-page", Type: &ReturnsStructRef{"Page"}},
		{DocURL: "https://developers.notion.com/reference/patch-page", Type: &ReturnsStructRef{"Page"}},
		{DocURL: "https://developers.notion.com/reference/retrieve-a-page-property", Type: &ReturnsInterface{"PropertyItemOrPropertyItemPagination"}},
		{DocURL: "https://developers.notion.com/reference/retrieve-a-block", Type: &ReturnsStructRef{"Block"}},
		{DocURL: "https://developers.notion.com/reference/update-a-block", Type: &ReturnsStructRef{"Block"}, Correct: func(method *MethodCoder) {
			for i, p := range method.Props.Doc.API.Params {
				if p.Name == "{type}" {
					if p.Desc != "The [block object `type`](ref:block#block-object-keys) value with the properties to be updated. Currently only `text` (for supported block types) and `checked` (for `to_do` blocks) fields can be updated." {
						panic(p.Desc)
					}

					text := p
					text.Name = "text"
					text.Type = "string" // TODO: rich text

					todo := p
					todo.Name = "to_do"
					todo.Type = "boolean" // TODO: *bool

					rest := append([]SSRPropsDocAPIParam{text, todo}, method.Props.Doc.API.Params[i+1:]...)
					method.Props.Doc.API.Params = append(method.Props.Doc.API.Params[0:i], rest...)
					break
				}
			}
		}},
		{DocURL: "https://developers.notion.com/reference/get-block-children", Type: &ReturnsStructRef{"BlockPagination"}},
		{DocURL: "https://developers.notion.com/reference/patch-block-children", Type: &ReturnsStructRef{"BlockPagination"}},
		{DocURL: "https://developers.notion.com/reference/delete-a-block", Type: &ReturnsStructRef{"Block"}},
		// https://developers.notion.com/reference/retrieve-a-comment
		// https://developers.notion.com/reference/create-a-comment
		// https://developers.notion.com/reference/get-user
		// https://developers.notion.com/reference/get-users
		// https://developers.notion.com/reference/get-self
	}

	builder := NewBuilder()
	eg := errgroup.Group{}

	for _, def := range methods {
		def := def
		coder := &MethodCoder{DocURL: def.DocURL, Type: def.Type}

		eg.Go(func() error {
			ssrProps, err := ParseMethod(coder.DocURL)
			if err != nil {
				return err
			}

			coder.Props = *ssrProps
			if def.Correct != nil {
				def.Correct(coder)
			}
			return nil
		})

		builder.Add(coder)
	}

	if err := eg.Wait(); err != nil {
		panic(err)
	}

	if err := builder.Build("../client.methods.go"); err != nil {
		panic(err)
	}
}
