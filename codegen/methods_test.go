package codegen

import (
	"testing"
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

	for _, method := range methods {
		method := method

		t.Run(method.DocURL, func(t *testing.T) {
			// t.Parallel()

			ssrProps, err := ParseMethod(method.DocURL)
			if err != nil {
				t.Fatal(err)
			}

			method.Props = *ssrProps
		})

		builder.Add(method)
	}

	if err := builder.Build("../client.methods.go"); err != nil {
		t.Fatal(err)
	}
}
