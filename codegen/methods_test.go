package codegen

import (
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestMethods(t *testing.T) {
	methods := []*MethodCoder{
		{DocURL: "https://developers.notion.com/reference/post-database-query", Returns: "PagePagination"},
		{DocURL: "https://developers.notion.com/reference/create-a-database", Returns: "Database"},
		{DocURL: "https://developers.notion.com/reference/update-a-database", Returns: "Database"},
		{DocURL: "https://developers.notion.com/reference/retrieve-a-database", Returns: "Database"},
		{DocURL: "https://developers.notion.com/reference/retrieve-a-page", Returns: "Page"},
		{DocURL: "https://developers.notion.com/reference/post-page", Returns: "Page"},
		{DocURL: "https://developers.notion.com/reference/patch-page", Returns: "Page"},
		{DocURL: "https://developers.notion.com/reference/retrieve-a-page-property", Returns: "PropertyItemOrPagination"},
		{DocURL: "https://developers.notion.com/reference/retrieve-a-block", Returns: "Block"},
		{DocURL: "https://developers.notion.com/reference/get-block-children", Returns: "BlockPagination"},
		{DocURL: "https://developers.notion.com/reference/patch-block-children", Returns: "BlockPagination"},
		{DocURL: "https://developers.notion.com/reference/delete-a-block", Returns: "Block"},
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
