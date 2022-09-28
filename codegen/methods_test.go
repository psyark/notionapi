package codegen

import (
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestMethods(t *testing.T) {
	methods := []*MethodCoder{
		{DocURL: "https://developers.notion.com/reference/post-database-query", Type: &ReturnsStructRef{"PagePagination"}},
		{DocURL: "https://developers.notion.com/reference/create-a-database", Type: &ReturnsStructRef{"Database"}},
		{DocURL: "https://developers.notion.com/reference/update-a-database", Type: &ReturnsStructRef{"Database"}},
		{DocURL: "https://developers.notion.com/reference/retrieve-a-database", Type: &ReturnsStructRef{"Database"}},
		{DocURL: "https://developers.notion.com/reference/retrieve-a-page", Type: &ReturnsStructRef{"Page"}},
		{DocURL: "https://developers.notion.com/reference/post-page", Type: &ReturnsStructRef{"Page"}},
		{DocURL: "https://developers.notion.com/reference/patch-page", Type: &ReturnsStructRef{"Page"}},
		{DocURL: "https://developers.notion.com/reference/retrieve-a-page-property", Type: &ReturnsInterface{"PropertyItemOrPropertyItemPagination"}},
		{DocURL: "https://developers.notion.com/reference/retrieve-a-block", Type: &ReturnsStructRef{"Block"}},
		{DocURL: "https://developers.notion.com/reference/get-block-children", Type: &ReturnsStructRef{"BlockPagination"}},
		{DocURL: "https://developers.notion.com/reference/patch-block-children", Type: &ReturnsStructRef{"BlockPagination"}},
		{DocURL: "https://developers.notion.com/reference/delete-a-block", Type: &ReturnsStructRef{"Block"}},
	}

	builder := NewBuilder()
	eg := errgroup.Group{}

	for _, method := range methods {
		method := method

		eg.Go(func() error {
			ssrProps, err := ParseMethod(method.DocURL)
			if err != nil {
				return err
			}

			method.Props = *ssrProps
			return nil
		})

		builder.Add(method)
	}

	if err := eg.Wait(); err != nil {
		panic(err)
	}

	if err := builder.Build("../client.methods.go"); err != nil {
		panic(err)
	}
}
